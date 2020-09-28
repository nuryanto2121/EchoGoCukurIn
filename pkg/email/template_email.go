package templateemail

const (
	VerifyCode = `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<title>Title of the document</title>
	</head>
	<body>

	<h4>Hai {Name}</h4>
	
	<h1>{GenerateCode}</h1>


	</body>
	</html>

	`

	SendPasswordCapster = `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<title>Register</title>
	</head>
	<body>

	<h4>Hai {Name}</h4>
	
	<h1>Password : {PasswordCode}</h1>

	Please change your password when first login.


	</body>
	</html>

	`
	SendRegister = `
	<!DOCTYPE html>
	<html lang="en">
	<head>
	<title>Informasi Login</title>
	</head>
	<body>

	<p>Hi {Name},</p>
	<p>Selamat, saat ini akun kamu telah aktif.</p>
	<p><strong>INFORMASI LOGIN</strong></p>
	<p>Username : {Email} <br/> PassWord : {PasswordCode}</p>
	<ul>
	<li>Usahakan agar kamu langsung mengganti password dan lakukan pergantian password secara berkala</li>
	<li>Password terdiri dari minimal 4 karakter &amp; maksimal 8 karakter dengan kombinasi huruf besar dan kecil, serta angka. Contoh : cuKuRin1</li>
	</ul>
	<p>&nbsp;</p>
	<p>KONTAK KAMI<br />
	Apabila ada pertanyaan, silahkan menghubungi kami melalui :</p>
	<ol>
	<li>Whatsapp 082308235470</li>
	<li>Email ke <a href="mailto:business@cukur-in.com">business@cukur-in.com</a></li>
	<li>Instagram @cukurin.id</li>
	</ol>
	<p>Terimakasih atas perhatian dan kepercayaannya.</p>
	<p>Salam</p>

	</body>
	</html>

	`
)
