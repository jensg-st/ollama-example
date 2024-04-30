package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	handler "github.com/jensg-st/ollama-example/service/pkg/handler"
)

func main() {
	r := chi.NewRouter()
	r.Post("/", handleRequest)
	http.ListenAndServe(":8080", r)
}

type request struct {
	Prompt string `json:"prompt"`
	Tuning string `json:"tuning"`
}

const (
	direktivActionIDHeader     = "Direktiv-ActionID"
	direktivErrorCodeHeader    = "Direktiv-ErrorCode"
	direktivErrorMessageHeader = "Direktiv-ErrorMessage"

	errCode = "io.direktiv.ollama.error"
	model   = "mistral"
)

func handleRequest(w http.ResponseWriter, r *http.Request) {

	aid := r.Header.Get(direktivActionIDHeader)
	log(aid, "got request")

	b, err := io.ReadAll(r.Body)
	if err != nil {
		respondWithErr(w, errCode, err.Error())
		return
	}
	defer r.Body.Close()

	var req request
	err = json.Unmarshal(b, &req)
	if err != nil {
		respondWithErr(w, errCode, err.Error())
		return
	}

	var t []byte
	if req.Tuning != "" {
		t, err = os.ReadFile(req.Tuning)
		if err != nil {
			respondWithErr(w, errCode, err.Error())
			return
		}
	}

	// add to prompt
	prompt := fmt.Sprintf("%s\n%s", string(t), req.Prompt)

	// request answer
	rh := handler.NewRequestHandler(model)
	err = rh.ProcessRequest(r.Context(), prompt)
	if err != nil {
		respondWithErr(w, errCode, err.Error())
		return
	}

	// prewpare response
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(rh.Response()))
}

func log(aid, l string) {
	if aid == "development" || aid == "Development" || aid == "" {
		fmt.Println(l)
	} else {
		http.Post(fmt.Sprintf("http://localhost:8889/log?aid=%s", aid), "plain/text", strings.NewReader(l))
	}
}

func respondWithErr(w http.ResponseWriter, code, err string) {
	w.Header().Set(direktivErrorCodeHeader, code)
	w.Header().Set(direktivErrorMessageHeader, err)
	w.WriteHeader(400)
	w.Write([]byte(err))
}

// GenerateResponseFunc func(GenerateResponse) error

// qwen
// '
// There are multiple automations available which can bestarted with a JSON file. The following is the list of automations and their JSON templates.

// Change Password:

// {
//     "action111": "change_password",
//     "user": "",
//     "server": "",
// }

// Start Server:

// {
//     "action": "start_server",
//     "server": "",
// }

// Stop Server:

// {
//     "action": "stop_server",
//     "server": "",
// }

// Which template would you pick and how would it look like if the user asks: Can you change the password for user tmpuser on server 192.168.1.10?
// Respond without text, only the valid JSON if all the information for the request is available. If not ask for the missing information.
// '

// '
// There are multiple automations available which can bestarted with a JSON file. The following is the list of automations and their JSON templates.

// Change Password:

// {
//     "action111": "change_password",
//     "user": "",
//     "server": "",
// }

// Start Server:

// {
//     "action": "start_server",
//     "server": "",
// }

// Stop Server:

// {
//     "action": "stop_server",
//     "server": "",
// }

// Which template would you pick and how would it look like if the user asks: I need to stop server 192.168.1.10?
// Respond without text, only the valid JSON if all the information for the request is available. If not ask for the missing information.
// '
