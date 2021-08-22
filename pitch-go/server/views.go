package server

import "fmt"

func RenderWesocketScript() string {
	return fmt.Sprintf(`
		<script>
		%s
			window.addEventListener('load', () => {
				w = new WebSocket('ws://' + window.location.host + window.location.pathname)

				w.addEventListener('message', (event) => {
					document.body.innerHTML = event.data;
				});
			});
		%s
		</script>
		`, "", "")
	// `, "/*", "*/")
}
