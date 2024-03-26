package posts_test

// func TestListPosts(t *testing.T) {
// 	creationDate, _ := time.Parse(customtime.DateFormat, "2000-01-01")

// 	postsStorage := repository.NewPosts(customtime.MockTimeProvider{}, &sync.Map{})
// 	userStorage := repository.NewUsers(customtime.MockTimeProvider{}, &sync.Map{})
// 	postsService := posts.NewPostsService(postsStorage, userStorage)

// 	author, _ := userStorage.GetUserByID(0)

// 	post1 := domain.Post{
// 		ID:          0,
// 		AuthorID:    0,
// 		Content:     "Заснял такие вот красивые деревья)",
// 		Attachments: []string{"tree1.jpeg", "tree2.jpeg", "tree3.jpeg"},
// 		CreatedAt:   customtime.CustomTime{Time: creationDate},
// 		UpdatedAt:   customtime.CustomTime{Time: creationDate},
// 	}

// 	post2 := domain.Post{
// 		ID:          1,
// 		AuthorID:    0,
// 		Content:     "Озеро недалеко от моего домика в Швейцарии. Красота!",
// 		Attachments: []string{"lake.jpeg"},
// 		CreatedAt:   customtime.CustomTime{Time: creationDate},
// 		UpdatedAt:   customtime.CustomTime{Time: creationDate},
// 	}

// 	post3 := domain.Post{
// 		ID:          2,
// 		AuthorID:    0,
// 		Content:     "Moя подруга - очень хороший фотограф",
// 		Attachments: []string{"camera.jpeg"},
// 		CreatedAt:   customtime.CustomTime{Time: creationDate},
// 		UpdatedAt:   customtime.CustomTime{Time: creationDate},
// 	}

// 	post4 := domain.Post{
// 		ID:          3,
// 		AuthorID:    0,
// 		Content:     "Мост в бесконечность",
// 		Attachments: []string{"bridge.jpeg"},
// 		CreatedAt:   customtime.CustomTime{Time: creationDate},
// 		UpdatedAt:   customtime.CustomTime{Time: creationDate},
// 	}

// 	post5 := domain.Post{
// 		ID:          4,
// 		AuthorID:    0,
// 		Content:     "Белые розы, белые розы... Не совсем белые, но все равно прекрасно)",
// 		Attachments: []string{"rose.jpeg"},
// 		CreatedAt:   customtime.CustomTime{Time: creationDate},
// 		UpdatedAt:   customtime.CustomTime{Time: creationDate},
// 	}

// 	tests := []struct {
// 		name string
// 		want []posts.PostWithAuthor
// 	}{
// 		{"List posts", []posts.PostWithAuthor{
// 			{Post: post1, Author: *author},
// 			{Post: post2, Author: *author},
// 			{Post: post3, Author: *author},
// 			{Post: post4, Author: *author},
// 			{Post: post5, Author: *author},
// 		}},
// 	}
// }
