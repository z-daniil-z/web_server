package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"web_server/dataBase"
)

func main() {
	config, err := getConfig("secureConfig.json")
	if err != nil {
		log.Fatalf(err.Error())
	}
	db, err := dataBase.NewDataBase(config[dataBaseField].(*dataBase.DBConfig))
	if err != nil {
		log.Fatalf(err.Error())
	}

	router := mux.NewRouter()
	handlers := NewHandlers(db, config[logField].(bool))
	staticDir := config[staticDirField].(string)

	router.HandleFunc("/registration", handlers.Registration).Methods("POST")
	router.HandleFunc("/login", handlers.Login).Methods("POST")
	router.HandleFunc("/profiles/{id}", handlers.GetUserInfo).Methods("GET")
	router.HandleFunc("/search", handlers.SearchUser).Methods("GET")
	router.HandleFunc("/delete", handlers.DeleteUser).Methods("POST")
	router.HandleFunc("/tags", handlers.GetTags).Methods("GET")
	router.HandleFunc("/tags/add", handlers.AddTagsToUser).Methods("POST")
	router.HandleFunc("/tags/get/{id}", handlers.GetUserTags).Methods("GET")
	router.HandleFunc("/tasks", handlers.GetTaskInfo).Methods("GET")
	router.HandleFunc("/tasks/tags", handlers.GetTaskTags).Methods("GET")
	router.HandleFunc("/invite/user/{id}", handlers.InviteUser).Methods("POST")
	router.HandleFunc("/invite/show", handlers.GetInvites).Methods("GET")
	router.HandleFunc("/validate/show", handlers.GetTasksToValidate).Methods("GET")
	router.HandleFunc("/quests/show", handlers.GetQuests).Methods("GET")
	router.HandleFunc("/quest/status/change", handlers.ChangeQuestStatus).Methods("POST")
	//For developers
	router.HandleFunc("/developers/getAccount", handlers.GetDeveloperAccount).Methods("GET")
	router.HandleFunc("/developers/postTag", handlers.PostTag).Methods("POST")
	router.HandleFunc("/developers/tasks/post", handlers.PostTask).Methods("POST")
	router.HandleFunc("/developers/tasks/addTags", handlers.AddTagsToTask).Methods("POST")

	go router.PathPrefix(staticDir).Handler(http.StripPrefix(staticDir, http.FileServer(http.Dir("."+staticDir))))

	go func() {
		if err := http.ListenAndServe(":80", http.HandlerFunc(handlers.redirectTLS)); err != nil {
			log.Fatalf("ListenAndServe error: %v", err)
		}
	}()

	err = http.ListenAndServeTLS(":443", config[certFilePathField].(string), config[keyFilePathField].(string), router)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)

	}
}
