package workers

func Init() {
	go syncConfig()
}
