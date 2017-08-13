package Global

type SerialOption struct {
	Portname string `json:"Portname"`
	Baudrate int `json:"Baudrate"`
	Databits int `json:"DataBits"`
	Stopbits int `json:"StopBits"`
	ParityMode int  `json:"ParityMode"`
	ReadInternval int  `json:"ReadInternval"` //Intervalo de tiempo (ms) que revisa si ha llegado nuevos datos
}
//Configuracion
type Configuration struct {
	SerialOption SerialOption `json:"SerialOption"` //Direccion Opciones de apertura del Puerto
	ServerTcpAddress string `json:"ServerTcpAddress"`
}

type ItemInfo struct {
	Value int	`json:"Value"`
	Date  string  `json:"TimeReceived"`
}
type Store struct {
	DataLast  map[string]*ItemInfo
	DataQueue map[string]*QueueItem
}

type Global struct {
	Config Configuration
	Store Store
}


//region QueueItem
type QueueItem []*ItemInfo

func (q *QueueItem) MaxItems() int {
	return 5;
}

func (q *QueueItem) Push(n *ItemInfo) {
	if q.Len()>=5 {
		q.Pop()
	}

	*q = append(*q, n)
}

func (q *QueueItem) Pop() (n *ItemInfo) {
	n = (*q)[0]
	*q = (*q)[1:]
	return
}

func (q *QueueItem) Len() int {
	return len(*q)
}

func (q QueueItem) ToList() []*ItemInfo {

	list := []*ItemInfo{}

	len := q.Len()

	for i:=0;i<len;i++ {
		list = append(list,q.Pop())
	}

	return list
}

//Obtiene el Primer elemento, pero sin sacarlo de la cola
func (q *QueueItem) GetLast() (n *ItemInfo){
	if (q.Len()>0) {
		n = (*q)[0]
	} else{
		n = nil
	}
	return n
}

func (q *QueueItem) GetFirst() (n *ItemInfo){
	if (q.Len()>0) {
		n = (*q)[q.Len()-1]
	} else{
		n = &ItemInfo{Value:0,Date:""}
	}
	return n
}

func (q *QueueItem) AllEqual(n int) bool {

	if (q.Len()<q.MaxItems()) {
		return false;
	}

	for i:=0;i<q.Len();i++ {
		if ((*q)[i]).Value != n {
			return false
		}
	}

	return true  //Todos los elementos son iguales
}

//endregion

var Resources Global

//var DB  map[string]*ItemInfo

func init() {
	Resources.Store.DataLast = make(map[string]*ItemInfo)
	Resources.Store.DataQueue = make(map[string]*QueueItem)

	Resources.Store.DataLast["A"]= &ItemInfo{0, ""}
	Resources.Store.DataLast["B"]= &ItemInfo{0, ""}
	Resources.Store.DataLast["C"]= &ItemInfo{0, ""}

	Resources.Store.DataQueue["A"]= &QueueItem{}
	Resources.Store.DataQueue["B"]= &QueueItem{}
	Resources.Store.DataQueue["C"]= &QueueItem{}

	//DB= make(map[string]*ItemInfo)
	//
	//DB["A"]=&ItemInfo{Value:5,Date: time.Now()}
	//DB["B"]=&ItemInfo{Value:5,Date: time.Now()}
	//DB["C"]=&ItemInfo{Value:5,Date: time.Now()}
}
