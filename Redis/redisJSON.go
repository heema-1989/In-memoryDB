package main

import "InMemoryDatabase/Redis/redis"

func main() {
	redis.SetKey("Name", "Heema", 0)
	redis.GetKey("Name")
}
