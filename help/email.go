package help

var EMAIL_PAYMENT_CANCEL = Parse(`Hi,

We just received subscription cancellation message from paypal.
You will be able to upload/download backups for 1 more month.
If you want to renew your subscription go to:

  https://baxx.dev/v1/sub/{{.PaymentID}}

and you will be redirected to paypal.com.


You can check the account status with:

  curl -u {{.Email}} -XPOST https://baxx.dev/protected/v1/status | json_pp

Thanks for using baxx.dev,
if you have any feedback please send me an email to jack@baxx.dev.

--
baxx.dev
`)

var EMAIL_AFTER_REGISTRATION = Parse(`Hi,
Thanks for registering to baxx.dev.

The servie I offer is still in Alpha stage, 
but I really appreciate the support.

## Plan (only one for now):

  Storage 10G
  Trial 1 Month 0.1E
  Subscription: 5E per Month
  Availability: ALPHA

  Here be Dragons! Data can be lost!

In order to use baxx.dev please subscribe at:

https://baxx.dev/v1/sub/{{.PaymentID}}
(this will redirect you to paypal)

## Tokens

 Tokens are like backup namespaces, you can have the same
 file in different tokens and it won't conflict.

Current Tokens:

{{ range .Tokens }}
  TOKEN: {{.UUID}}
  Write Only: {{ .WriteOnly }}
  Keep N Versions {{ .NumberOfArchives }}

{{end}}
* Create New Tokens:

 curl -u {{ .Email }} -d '{"write_only":false, "keep_n_versions":7}' \
   https://baxx.dev/protected/v1/create/token

 + Write Only:
   tokens can only add but not get files (without password)

 + Keep #N Versions:
   How many versions per file (with different sha256) to keep.
   Useful for database or modified files archives like, e.g:

    mysqldump | curl curl --data-binary @- \
     https://baxx.dev/v1/io/$TOKEN/mysql.gz

* Delete tokens

 curl -u {{ .Email }} -d '{"uuid": "TOKEN-UUID"}' \
   https://baxx.dev/protected/v1/delete/token

 this will delete the token and all the files in it

## File operations

* File Upload:

 cat path/to/file | curl --data-binary @- \
   https://baxx.dev/v1/io/$TOKEN/path/to/file

 Same filepath can have up to #N Versions depending on the token
 configuration.

 Uploading the same sha256 resulting in reusing existing version
 and also does not consume quota.

* File Download:

 curl https://baxx.dev/v1/io/$TOKEN/path/to/file > file

 Downloads the last upload version

* File Delete:

 curl -XDELETE https://baxx.dev/v1/io/$TOKEN/path/to/file

 deletes all versions of a file

* List Files in path LIKE /path/to%:

 curl https://baxx.dev/v1/ls/$TOKEN/path/to

 use -H "Accept: application/json" if you want json back
 by default it prints human readable text


WriteOnly tokens require BasicAuth and /protected prefix.

* Download from WriteOnly token:

 curl -u {{ .Email }} \
   https://baxx.dev/protected/v1/io/$TOKEN/path/to/file

* Delete with WriteOnly token:

 curl -u {{ .Email }} -XDELETE \
   https://baxx.dev/v1/io/$TOKEN/path/to/file

* List with WriteOnly token:

 curl -u {{ .Email }} \
   https://baxx.dev/protected/v1/ls/$TOKEN/path/to/


## Profile Management

* Register:

 curl -d '{"email":"{{.Email}}", "password":"mickey mouse"}' \
  https://baxx.dev/v1/register | json_pp

* Change Password

 curl -u {{.Email}} -d'{"new_password": "donald mouse"}' \
  https://baxx.dev/protected/v1/replace/password | json_pp

 (use https://www.xkcd.com/936/)

* Change Email

 curl -u {{.Email}} -d'{"new_email": "x@example.com"}' \
  https://baxx.dev/protected/v1/replace/email | json_pp

 It will also send new verification email, you can
 also use the replace/email endpoint to resend the
 verification email.

* User Status

 curl -u {{.Email}} -XPOST https://baxx.dev/protected/v1/status

 shows things like
  * is the email verified
  * is subscription active [ not done yet ]
  * current tokens
  * size used

--
baxx.dev
`)

var EMAIL_VALIDATION = Parse(`Hi,

this is the verification link: 

  https://baxx.dev/v1/verify/{{.ID}}

You can check the account status with:

  curl -u {{.Email}} -XPOST https://baxx.dev/protected/v1/status | json_pp

PS:
It is very likely that this email goes to the spam folder 
because it is small and texty.. anyway, I hope it doesnt.
(fingers crossed)

--
baxx.dev
`)

var EMAIL_PAYMENT_THANKS = Parse(`Hi,

Thanks for subscribing!
Even though the service is just in alpha state, it is much
appreciated!

If you want to cancel you have to do that in your paypal account.


You can check the account status with:

  curl -u {{.Email}} -XPOST https://baxx.dev/protected/v1/status | json_pp

--
baxx.dev
`)