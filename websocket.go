package dewu

import "github.com/gorilla/websocket"

type WebSocket struct {
	Pixelclick PixelClick      `json:"pixelclick"`
	WsConn     *websocket.Conn `json:"ws_conn"`
}
type PixelClick struct {
	*Pixel
	Lastclick int `json:"lastclick"`
}
