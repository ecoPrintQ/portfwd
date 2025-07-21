# Troubleshooting 

## 1- Puerto en uso 

Si el sistema genera el siguiente log:
```
"Error. The connection failed (ListenPaket): listen udp4 0.0.0.0:161: bind: An attempt was made to access a socket in a way forbidden by its access permissions."
```
Verifique que el puerto que intenta usar se encuentre disponible. 
Puede usar el comando (reemplace PORT por el puerto a consultar):
```bash
    netstat -ano | findstr :PORT
```

Identifique el proceso que hace uso del puerto (reemplace PID_PORT por el valor obtenido arriba):
```bash
tasklist /FI "PID eq PID_PORT"
```
Debe detener el servicio que está haciendo uso del puerto o cambiar el puerto local que recibe la solicitud.