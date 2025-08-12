package ui

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"tcllike/evaluator"
	"tcllike/lexer"
	"tcllike/parser"
)

type WebUI struct {
	tmpl *template.Template
}

func NewWebUI() *WebUI {
	tmpl := template.Must(template.New("web").Parse(`
	<!DOCTYPE html>
	<html>
	<head>
		<title>Tcl-like Web REPL</title>
		<style>
			body { font-family: monospace; margin: 20px; }
			#repl { border: 1px solid #ccc; padding: 10px; height: 400px; overflow-y: scroll; }
			#input { width: 80%; padding: 5px; }
			button { padding: 5px 15px; }
			.prompt { color: #0099cc; font-weight: bold; }
			.result { color: #00cc99; }
			.error { color: #ff3333; font-weight: bold; }
		</style>
	</head>
	<body>
		<h1>Tcl-like Web REPL</h1>
		<div id="repl"></div>
		<form id="form">
			<input id="input" type="text" placeholder="Enter Tcl command (e.g. puts hello)" autofocus>
			<button type="submit">Eval</button>
		</form>
		<script>
			const repl = document.getElementById('repl');
			const form = document.getElementById('form');
			const input = document.getElementById('input');
			
			form.addEventListener('submit', async (e) => {
				e.preventDefault();
				const code = input.value;
				if (!code) return;
				
				const prompt = document.createElement('div');
				prompt.className = 'prompt';
				prompt.textContent = 'tcl> ' + code;
				repl.appendChild(prompt);
				
				const response = await fetch('/eval', {
					method: 'POST',
					headers: { 'Content-Type': 'application/json' },
					body: JSON.stringify({ code })
				});
				
				const result = await response.json();
				const output = document.createElement('div');
				
				if (result.error) {
					output.className = 'error';
					output.textContent = '✗ ' + result.error;
				} else {
					output.className = 'result';
					output.textContent = '=> ' + result.output;
				}
				
				repl.appendChild(output);
				input.value = '';
				repl.scrollTop = repl.scrollHeight;
			});
		</script>
	</body>
	</html>
	`))
	return &WebUI{tmpl: tmpl}
}

func (w *WebUI) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/":
		w.tmpl.Execute(res, nil)
	case "/eval":
		var data struct {
			Code string `json:"code"`
		}
		if err := json.NewDecoder(req.Body).Decode(&data); err != nil {
			http.Error(res, err.Error(), http.StatusBadRequest)
			return
		}

		ast, err := parser.Parse(lexer.Tokenize(data.Code))
		if err != nil {
			json.NewEncoder(res).Encode(map[string]string{"error": err.Error()})
			return
		}

		result, err := evaluator.Eval(ast)
		if err != nil {
			json.NewEncoder(res).Encode(map[string]string{"error": err.Error()})
			return
		}

		json.NewEncoder(res).Encode(map[string]string{"output": fmt.Sprintf("%v", result)})
	default:
		http.NotFound(res, req)
	}
}
