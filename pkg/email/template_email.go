package templateemail

const (
	VerifyCode = `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<title>Forgot Password</title>
	</head>
	<body>

	<p>Hallo {Name},</p>
	<p style="text-align: justify;">Anda telah melakukan permohonan lupa password untuk Akun Anda. Abaikan Email ini Jika Anda tidak merasa melakukan permintaan lupa password.</p>
	<p style="text-align: justify;"><br />Silahkan masukan kode OTP di bawah ini untuk memasukkan password baru Anda. OTP ini hanya berlaku 1x24 jam sejak Email ini dikirimkan atau setelah password Anda berhasil diubah.</p>
	<p style="text-align: justify;"></p>
	<table style="height: 35px; width: 93.0782%; border-collapse: collapse; margin-left: auto; margin-right: auto;" height="22" border="1">
	<tbody>
		<tr>
			<td style="width: 100%; background-color: #00ffff; text-align: center;">
				<h4><span style="color: #0000ff;">{OTP}</span></h4>
			</td>
		</tr>
	</tbody>
	</table>
	<p></p>
	<p></p>



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
