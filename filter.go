package rxgo

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"time"
)

// 过滤操作结构体
type filterOp struct {
	opFunc func(ctx context.Context, o *Observable, item reflect.Value, out chan interface{}) (end bool)
}

// op将前一个被观察者的输入流当作本被观察者的输入
func (tsop filterOp) op(ctx context.Context, o *Observable) {
	in := o.pred.outflow
	out := o.outflow
	var wg sync.WaitGroup

	go func() {
		end := false
		for x := range in {
			if end {
				continue
			}
			xv := reflect.ValueOf(x)
			if e, ok := x.(error); ok && !o.flip_accept_error {
				o.sendToFlow(ctx, e, out)
				continue
			}
			switch threading := o.threading; threading {
			case ThreadingDefault:
				if tsop.opFunc(ctx, o, xv, out) {
					end = true
				}
			case ThreadingIO:
				fallthrough
			case ThreadingComputing:
				wg.Add(1)
				go func() {
					defer wg.Done()
					if tsop.opFunc(ctx, o, xv, out) {
						end = true
					}
				}()
			default:
			}
		}

		if o.flip != nil {
			buf := reflect.ValueOf(o.flip)
			if buf.Kind() != reflect.Slice {
				panic("flip is not a slice")
			}
			for i:=0; i<buf.Len(); i++ {
				o.sendToFlow(ctx, buf.Index(i).Interface(), out)
			}
		}

		// 等待所有go程结束并关闭outflow
		wg.Wait()
		o.closeFlow(out)
	}()
}

// Debounce间隔timespan时间输出item；
func (parent *Observable) Debounce(timespan time.Duration) (o *Observable){
	o = parent.newTransformObservable("Debounce")
	o.threading = ThreadingComputing
	count := 0
	o.operator = filterOp{
		func(ctx context.Context, o *Observable, item reflect.Value, out chan interface{}) (end bool){
			count++
			var tempCount = count
			time.Sleep(timespan)
			if tempCount == count{
				end = o.sendToFlow(ctx, item.Interface(), out)
			}
			return
		}}
	return o
}

// Distinct去除输入流中的重复item
func (parent *Observable) Distinct()(o *Observable){
	o = parent.newTransformObservable("Distinct")
	var uni = map[string]bool{}
	o.operator = filterOp{
		func(ctx context.Context, o *Observable, item reflect.Value, out chan interface{}) (end bool){
			itemStr := fmt.Sprintf("%v", item)
			if _, ok := uni[itemStr]; !ok{
				uni[itemStr] = true
				end = o.sendToFlow(ctx, item.Interface(), out)
			}
			return false
		}}
	return o
}

// ElementAt 返回第id位的元素
func (parent *Observable) ElementAt (id int)(o *Observable){
	o = parent.newTransformObservable("ElementAt")
	count := 0
	o.operator = filterOp{
		func(ctx context.Context, o *Observable, item reflect.Value, out chan interface{}) (end bool){
			if count == id{
				end = o.sendToFlow(ctx, item.Interface(), out)
			}
			count++
			return
		}}
	return o
}

// First发射首位元素
func (parent *Observable) First ()(o *Observable){
	o = parent.newTransformObservable("First")
	o.operator = filterOp{
		func(ctx context.Context, o *Observable, item reflect.Value, out chan interface{}) (end bool){
			o.sendToFlow(ctx, item.Interface(), out)
			return true
		}}
	return o
}

// IgnoreElements忽略输入流中所有item，只发送报错结束/结束信息；
func (parent *Observable) IgnoreElements ()(o *Observable){
	o = parent.newTransformObservable("IgnoreElements")
	o.operator = filterOp{
		func(ctx context.Context, o *Observable, item reflect.Value, out chan interface{}) (end bool){
			return
		}}
	return o
}

// last发射末位元素
func (parent *Observable) Last()(o *Observable){
	o = parent.newTransformObservable("last")
	o.operator = filterOp{
		func(ctx context.Context, o *Observable, item reflect.Value, out chan interface{}) (end bool){
			o.flip = append([]interface{}{}, item.Interface())
			return false
		}}
	return o
}


// Sample定期发射被观察者最近发射的数据项
func (parent *Observable) Sample(timespan time.Duration)(o *Observable){
	o = parent.newTransformObservable("Sample")
	o.threading = ThreadingComputing    //2
	temp := make([]reflect.Value, 256)
	timer := false
	var wa sync.WaitGroup
	count := 0
	span := + 50*time.Millisecond
	o.operator = filterOp{
		opFunc:func(ctx context.Context, o *Observable, item reflect.Value, out chan interface{}) (end bool){
			temp[count] = item
			count++
			if !timer{
				wa.Add(1)
				go func() {
					defer wa.Done()
					for {
						time.Sleep(timespan + span)
						select{
						case <-ctx.Done():
							return
						default:
							if count == 0 {
								return
							}
							if o.sendToFlow(ctx, temp[count-1].Interface(), out){
								return
							}
							count = 0
						}
					}
				}()
				timer = true
				wa.Wait()
			}
			return false

		}}
	return o
}

// skip跳过前n项
func (parent *Observable) Skip (n int)(o *Observable){
	o = parent.newTransformObservable("Skip")
	index := 0
	o.operator = filterOp{
		opFunc:func(ctx context.Context, o *Observable, item reflect.Value, out chan interface{}) (end bool){
			if index >= n{
				end = o.sendToFlow(ctx, item.Interface(), out)
			}
			index++
			return
		}}
	return o
}

// Skiplast跳过后n项发送
func (parent *Observable) Skiplast (n int)(o *Observable){
	o = parent.newTransformObservable("Skiplast")
	var temp []reflect.Value
	index := 0
	o.operator = filterOp{
		opFunc:func(ctx context.Context, o *Observable, item reflect.Value, out chan interface{}) (end bool){
			if index >= n{
				end = o.sendToFlow(ctx, temp[0].Interface(), out)
				temp = temp[1:]
			}
			index++
			temp = append(temp, item)
			return
		}}
	return o
}

// Take只发射前n项
func (parent *Observable) Take (n int)(o *Observable){
	o = parent.newTransformObservable("Take")
	var index = 0
	o.operator = filterOp{
		func(ctx context.Context, o *Observable, item reflect.Value, out chan interface{}) (end bool){
			if index < n{
				end = o.sendToFlow(ctx, item.Interface(), out)
			}
			index++
			return
		}}
	return o
}





