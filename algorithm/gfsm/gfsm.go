package gfsm

import (
	"container/list"
	"fmt"
	"reflect"
)

type callBackFunc func(arg interface{})

// FSMState State父struct
type FSMState struct {
	id       string                    //状态对应的键值
	callback callBackFunc              //普通的回调函数
	funName  string                    //对象对调方法的名字
	object   interface{}               //创建状态机的对象
	methods  map[string]reflect.Method //创装状态机对象的所有方法的集合
	arg      interface{}               //回调方法的参数
}

/*
创建状态机的状态方法
@param1 id:      状态对应的id
@param2: fun:    普通函数回调函数,兼容接口
@param3: arg:    回调函数的参数
@param4: object: 创建状态机的对象
@param5: name:   对象的回调方法，兼容接口
注意：两个兼容的回调函数，区别是：一个是函数，一个是创建状态机的对象的方法
*/
func CreateMatchineState(id string, fun callBackFunc, arg interface{}, object interface{}, name string) *FSMState {

	methods := make(map[string]reflect.Method)
	if object != nil {
		typ := reflect.TypeOf(object)
		for m := 0; m < typ.NumMethod(); m++ {
			method := typ.Method(m)
			mname := method.Name
			methods[mname] = method
		}
	}

	return &FSMState{
		id:       id,
		callback: fun,
		arg:      arg,
		object:   object,
		methods:  methods,
		funName:  name,
	}
}

/*
状态方法：进入状态
*/
func (fs *FSMState) Enter() {

}

/*
状态方法：状态处理函数
*/
func (fs *FSMState) Do() {
	if fs.object != nil {

		params := make([]reflect.Value, 2)
		params[0] = reflect.ValueOf(fs.object)
		params[1] = reflect.ValueOf(fs.arg)
		fs.methods[fs.funName].Func.Call(params)
		return
	}
	fs.callback(fs.arg)
}

/*
状态方法：退出状态
*/
func (fs *FSMState) Exit() {
}

/*
状态方法：添加状态回调函数
*/
func (fs *FSMState) addStateCallBack(f callBackFunc) {
	fs.callback = f
}

/*
状态方法：状态转移检测
*/
func (fs *FSMState) CheckTransition() {
	//
}

/*
设置回调参数值
*/
func (fs *FSMState) setCallBackArg(arg interface{}) {

	fs.arg = arg
}

/*******************************************************************************************/
func CreateFSM() *FSM {
	it := &FSM{}
	it.Init()
	return it
}

type FSM struct {
	// 持有状态集合
	statesMap map[string]*FSMState
	//
	action *list.List
	// 当前状态
	current_state *FSMState
	// 下一个状态
	next_state *FSMState
	// 默认状态
	default_state *FSMState
	runState      int
}

/*
状态机方法：初始化FSM
*/
func (f *FSM) Init() {
	//
	f.runState = 0
	f.statesMap = make(map[string]*FSMState)
	f.action = list.New()
}

/*
状态机方法：启动状态机
*/
func (f *FSM) Start() {
	if f.action.Len() == 0 {
		fmt.Println("FSM-Start():leave1")
		return
	}
	firstKey := f.action.Front()
	f.current_state = f.statesMap[firstKey.Value.(string)]
	var index string = f.CalcNextStateKey(f.current_state)
	f.next_state = f.statesMap[index]

	f.runState = 1
	f.DoFsmState()
}

/*
设状态机方法：置默认的State
@param: state:状态
*/
func (f *FSM) SetDefaultState(state *FSMState) {
	f.default_state = state
}

/*
添状态机方法：加状态到FSM
@param1 : 状态对应的key
@param : state:状态
*/
func (f *FSM) AddState(key string, state *FSMState) {
	f.statesMap[key] = state
	f.action.PushBack(key)
}

/*
根状态机方法：据key获取状态
@param: key:状态的Id
*/
func (f *FSM) GetStateById(key string) *FSMState {
	return f.statesMap[key]
}

/*
获状态机方法：取当前状态
*/
func (f *FSM) GetCurrentState() *FSMState {
	return f.current_state
}

/*
获状态机方法：取下一下状态
*/
func (f *FSM) GetNextState() *FSMState {
	return f.next_state
}

/*
状状态机方法：态切换
*/
func (f *FSM) DoFsmState() {
	if f.runState != 1 {
		return
	}
	//执行当前状态的动作
	f.current_state.Do()
}

func (f *FSM) SwitchFsmState() {
	f.current_state = f.next_state
	var index string = f.CalcNextStateKey(f.current_state)
	f.next_state = f.statesMap[index]
	f.DoFsmState()
}

/*
状态机方法：计算下一个状态
@parma1: crrrent:当前的状态
*/
func (f *FSM) CalcNextStateKey(current *FSMState) string {
	var index string = f.default_state.id
	for e := f.action.Front(); e != nil; e = e.Next() {
		if e.Value == current.id {
			if e.Next() != nil {
				index = (e.Next().Value).(string)
			}
			break
		}
	}
	return index
}

/*
状态机方法：设置状态的执行参数
@param1: key:状态对应的Id
@param2: arg:状态的回调函数执行参数
*/
func (f *FSM) SetStateCallBackArg(key string, arg interface{}) {
	state := f.GetStateById(key)
	if state == nil {
		return
	}
	state.setCallBackArg(arg)
}

/*
暂状态机方法：停FSM
*/
func (f *FSM) PauseStateMachine() {
	f.runState = 2
}

/*
重状态机方法：置FSM
*/
func (f *FSM) ResetStateMachine() {
	f.runState = 1
	f.current_state = f.default_state
	// 下一个状态
	var index string = f.CalcNextStateKey(f.current_state)
	f.next_state = f.statesMap[index]
}
