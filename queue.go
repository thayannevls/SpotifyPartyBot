package main

type Queue struct {
	users   []*User
	current int
	playing bool
}

func (queue *Queue) Add(user *User) {
	queue.users = append(queue.users, user)
}

func (queue *Queue) Remove(userID string) {
	for index, user := range queue.users {
		if user.discord.ID == userID {
			queue.users = append(queue.users[:index], queue.users[index+1:]...)
			return
		}
	}
}
