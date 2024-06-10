package utils

func EmailStyle(name string, company string) string {

    emailBody := `
    <html>
    <head>
        <style>
            body {
                font-family: Arial, sans-serif;
                margin: 0;
                padding: 0;
                background-color: #f4f4f4;
            }
            .container {
                width: 80%;
                margin: 0 auto;
                padding: 20px;
                background-color: #ffffff;
                box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
            }
            .header {
                text-align: center;
                padding: 10px 0;
            }
            .header img {
                width: 120px;
                height: auto;
            }
            .content {
                text-align: left;
                padding: 20px;
            }
            .content h2 {
                color: ##067945;
            }
            .content p {
                color: #26a345;
                font-size: 16px;
                line-height: 1.5;
            }
        </style>
    </head>
    <body>
        <div class="container">
            <div class="header">
                <img src="https://assets-global.website-files.com/5c995cb29175d63ad9e6ba7a/5f595a6c6952853b021e4ce6_Forem%20Icon%203-01.png" alt="Forem Logo">
                <br>
				<h2>Hola ` + name + `</h2>
            </div>
			<br>
			<br>
            <div class="content">
                <p>La búsqueda de información de ` + company + ` ha finalizado. A continuación, se encuentra el archivo adjunto.</p>
            </div>
			<br>
			<br>
        </div>
    </body>
    </html>
    `
    return emailBody
}