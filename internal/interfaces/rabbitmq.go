package interfaces

type IRbmqWorker interface {
	HandleLogFromQueue(msg []byte) error
	HandleEmailFromQueue(msg []byte) error
}
