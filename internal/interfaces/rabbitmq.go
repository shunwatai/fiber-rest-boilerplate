package interfaces

type IRbmqWorker interface {
	HandleLogFromQueue(msg []byte)
	HandleEmailFromQueue(msg []byte)
}
