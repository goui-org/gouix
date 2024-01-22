const ws = new WebSocket('ws://' + window.location.host + '/hot')
ws.onmessage = e => {
	if (e.data === 'reload') return window.location.reload()
	const overlay = document.createElement('div')
	overlay.style = 'position: fixed; left: 0; right: 0; top: 0; bottom: 0; background: #000c; color: #e77; font-size: 18px'
	const msg = document.createElement('div')
	msg.style = 'width: 100%; max-width: 600px; margin: auto; margin-top: 5vh; line-height: 200%; padding: 0 20px;'
	msg.innerText = e.data
	overlay.appendChild(msg)
	document.body.appendChild(overlay)
}
ws.onopen = () => ws.send('loaded')
ws.onclose = () => window.close()
