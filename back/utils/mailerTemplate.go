package utils

func EmailStyle (name string, company string, data string) string {

	emailBody := "<html><h2>Hola " + name  + "</br><p>A continuacion esta la informacion obtenida de " + company + "</p><br>" + data + "</html>"

	return emailBody
}