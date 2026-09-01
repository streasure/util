package msg

type Message interface {
	Handle()
	Wait()
	Done()
}

type SyncMessage struct {
	done chan struct{}
}

func NewSyncMessage() *SyncMessage {
	return &SyncMessage{
		done: make(chan struct{}),
	}
}

func (m *SyncMessage) Done() {
	close(m.done)
}

func (m *SyncMessage) Wait() {
	<-m.done
}

type AsyncMessage struct{}

func NewAsyncMessage() *AsyncMessage {
	return new(AsyncMessage)
}

func (m AsyncMessage) Done() {}

func (m AsyncMessage) Wait() {}
