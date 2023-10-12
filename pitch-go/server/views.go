package server

func RenderWesocketScript() string {
	return `
		<script>
			window.addEventListener('load', () => {
				const w = new WebSocket('ws://' + window.location.host + window.location.pathname)

				w.addEventListener('message', (event) => {
					document.body.innerHTML = event.data;
				});
			});
		</script>
	`
}
