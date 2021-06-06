package main

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

func conexionBD() (conexion *sql.DB) {
	Driver := "mysql"
	Usuario := "root"
	Contraseña := "password"
	NombreBase := "sistema"

	conexion, err := sql.Open(Driver, Usuario+":"+Contraseña+"@tcp(127.0.0.1)/"+NombreBase)
	if err != nil {
		panic(err.Error())
	}

	return conexion
}

var plantillas = template.Must(template.ParseGlob("plantillas/*"))

func main() {
	http.HandleFunc("/", inicio)
	http.HandleFunc("/crear", crear)
	http.HandleFunc("/insertar", insertar)
	http.HandleFunc("/borrar", borrar)
	http.HandleFunc("/editar", editar)
	http.HandleFunc("/actualizar", actualizar)
	log.Println("Servidor corriendo ...")
	http.ListenAndServe(":8080", nil)
}

type Empleado struct {
	Id     int
	Nombre string
	Correo string
}

func inicio(w http.ResponseWriter, r *http.Request) {

	conexionEstablecida := conexionBD()
	registros, err := conexionEstablecida.Query("select * from empleados")
	if err != nil {
		panic(err.Error())
	}

	empleado := Empleado{}
	arregloEmpleado := []Empleado{}

	for registros.Next() {
		var id int
		var nombre string
		var correo string
		err = registros.Scan(&id, &nombre, &correo)
		if err != nil {
			panic(err.Error())
		}
		empleado.Id = id
		empleado.Nombre = nombre
		empleado.Correo = correo
		arregloEmpleado = append(arregloEmpleado, empleado)

	}

	plantillas.ExecuteTemplate(w, "inicio", arregloEmpleado)
}

func crear(w http.ResponseWriter, r *http.Request) {
	plantillas.ExecuteTemplate(w, "crear", nil)

}

func insertar(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		nombre := r.FormValue("nombre")
		correo := r.FormValue("correo")

		conexionEstablecida := conexionBD()
		insertarRegistros, err := conexionEstablecida.Prepare("Insert into empleados (nombre,correo) value (?,?)")
		if err != nil {
			panic(err.Error())
		}

		insertarRegistros.Exec(nombre, correo)
		code := 301
		http.Redirect(w, r, "/", code)

	}
}

func borrar(w http.ResponseWriter, r *http.Request) {

	idEmpleado := r.URL.Query().Get("id")
	conexionEstablecida := conexionBD()
	borrarRegistro, err := conexionEstablecida.Prepare("Delete from empleados where id =?")

	if err != nil {
		panic(err.Error())
	}
	borrarRegistro.Exec(idEmpleado)
	code := 301
	http.Redirect(w, r, "/", code)

}

func editar(w http.ResponseWriter, r *http.Request) {
	idEmpleado := r.URL.Query().Get("id")
	conexionEstablecida := conexionBD()
	editarRegistro, err := conexionEstablecida.Query("Select * from empleados where id =?", idEmpleado)
	if err != nil {
		panic(err.Error())
	}

	empleado := Empleado{}

	for editarRegistro.Next() {
		var id int
		var nombre, correo string
		err = editarRegistro.Scan(&id, &nombre, &correo)
		if err != nil {
			panic(err.Error())

		}

		empleado.Id = id
		empleado.Nombre = nombre
		empleado.Correo = correo

	}

	plantillas.ExecuteTemplate(w, "editar", empleado)

}

func actualizar(w http.ResponseWriter, r *http.Request) {

	if r.Method == "POST" {
		id := r.FormValue("id")
		nombre := r.FormValue("nombre")
		correo := r.FormValue("correo")

		conexionEstablecida := conexionBD()

		actualzarRegistro, err := conexionEstablecida.Prepare("Update empleados set nombre=? , correo=?  where id =?")
		if err != nil {
			panic(err.Error())
		}
		actualzarRegistro.Exec(nombre, correo, id)
		code := 301
		http.Redirect(w, r, "/", code)
	}
}
