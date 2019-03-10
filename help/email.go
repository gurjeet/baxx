package help

var EMAIL_PAYMENT_CANCEL = Parse(`Hi,

We just received subscription cancellation message from paypal.
You will be able to upload/download backups for 1 more month.
If you want to renew your subscription go to:

https://baxx.dev/sub/{{.PaymentID}}

and you will be redirected to paypal.com.

You can check the account status with:

curl -u {{.Email}} -XPOST https://baxx.dev/protected/status | json_pp

Thanks for using baxx.dev,
if you have any feedback please send me an email to jack@baxx.dev.

--
baxx.dev
`)

var EMAIL_AFTER_REGISTRATION = Parse(`Hi,

The service I offer is still in Alpha stage, but I really appreciate
the support.

# Subscription

## Plan (only one for now):

Storage 10G
Trial 1 Month 0.1E
Subscription: 5E per Month
Availability: ALPHA

Here be Dragons! Data can be lost!

## Subscribe

In order to use baxx.dev you need a subscription,
At the moment I support only paypal.com, please visit:

https://baxx.dev/sub/{{.PaymentID}}
To be redirected to paypal.com.

## Verify your email

Email verification is also required, you should've received the
verification link in another email.
{{ if .LastVerificationID }}
Or you could also click on:

https://baxx.dev/verify/{{.LastVerificationID}}

{{ end }}
Thanks again!

# Tokens

Tokens are like backup namespaces, you can have the same file in
different tokens and it won't conflict.

There are 2 kinds of tokens, ReadWrite and WriteOnly,
ReadWrite tokens dont require any credentials for create, delete and
list files, WriteOnly tokens require credentials for *list* and
*delete*.

## Current Tokens:

{{ range .Tokens }}
  TOKEN: {{.UUID}}
    {{ if .Name }}Name: {{ .Name }}{{ end }}
    Write Only: {{ .WriteOnly }}
    Keep N Versions {{ .NumberOfArchives }}
{{end}}
## Create New Tokens:

curl -u {{ .Email }} \
 -d '{"write_only":false, "keep_n_versions":7, "name": "example"}' \
 https://baxx.dev/protected/create/token

Write Only:
 tokens can only add but not get files (without password)

Keep #N Versions:
 How many versions per file (with different sha256) to keep.  Useful
 for database or modified files archives like, e.g:

 mysqldump | curl --data-binary @- \
  https://baxx.dev/io/$TOKEN/mysql.gz
## Modify tokens

curl -u {{ .Email }} \
 -d '{"write_only":false,token:"TOKEN-UUID","name":"example"}' \
 https://baxx.dev/protected/change/token
## Delete tokens

curl -u {{ .Email }} -d '{"token": "TOKEN-UUID"}' \
 https://baxx.dev/protected/delete/token

this will delete the token and all the files in it

# File operations

## File Upload:

cat path/to/file | curl --data-binary @- \
 https://baxx.dev/io/$TOKEN/path/to/file

curl -T path/to/file https://baxx.dev/io/$TOKEN/path/to/file

Same filepath can have up to #N Versions depending on the token
configuration.

Uploading the same sha256 resulting in reusing existing version and
also does not consume quota.

## File Download:

curl https://baxx.dev/io/$TOKEN/path/to/file > file

Downloads the last upload version

## File Delete:

Delete single file:
curl -XDELETE https://baxx.dev/io/$TOKEN/path/to/file

Delete all files in a directory, but not the subdirectories:
curl -d '{"force":true}' https://baxx.dev/io/$TOKEN/path

## List Files in path LIKE /path/to%:

curl https://baxx.dev/ls/$TOKEN/path/to

use -H "Accept: application/json" if you want json back by default it
prints human readable text


## WriteOnly Tokens

Write Only tokens require BasicAuth.
The idea is that you can put them in in-secure places and not worry
about someone reading your data if they get stolen.

## Download from WriteOnly token:

curl -u {{ .Email }} \
 https://baxx.dev/io/$TOKEN/path/to/file

## Delete with WriteOnly token:

curl -u {{ .Email }} -XDELETE \
 https://baxx.dev/io/$TOKEN/path/to/file

## List with WriteOnly token:

curl -u {{ .Email }} \
 https://baxx.dev/ls/$TOKEN/path/to/


# Profile Management

## Register:

curl -d '{"email":"{{.Email}}", "password":"mickey mouse"}' \
 https://baxx.dev/register | json_pp

## Change Password

curl -u {{.Email}} -d'{"new_password": "donald mouse"}' \
 https://baxx.dev/protected/replace/password | json_pp

(use https://www.xkcd.com/936/)

## Change Email

curl -u {{.Email}} -d'{"new_email": "x@example.com"}' \
 https://baxx.dev/protected/replace/email | json_pp

It will also send new verification email, you can also use the
replace/email endpoint to resend the verification email.

## User Status

curl -u {{.Email}} -XPOST https://baxx.dev/protected/status

shows things like
 * is the email verified
 * is subscription active [ not done yet ]
 * current tokens
 * size used

# Encryption

Your data is compressed and encrypted when received, the encryption
key is auto generated uuid, and the purpose of the encryption is
simply to obscure the data in case the machines are hacked, hacker
will have to also get access to the database as well.

Anyway, dont trust it and use encryption when uploading.

# Examples

## upload everything from a directory

find . -type f -exec curl --data-binary @{}      \
              https://baxx.dev/io/$BAXX_TOKEN/{} \;

## upload only the files that have difference in shasum

for i in $(find . -type f); do \
 echo -n "$i.."
 sha=$(shasum -a 256 $i | cut -f 1 -d ' ')
 (curl -s https://baxx.dev/sha256/$BAXX_TOKEN/$sha -f && echo SKIP $i) || \
 (curl -T $i https://baxx.dev/io/$BAXX_TOKEN/$i -f)
done


## shell alias
# indentation is messed up  to fit 80 chars

---

export BAXX_TOKEN=...
baxx_put() {
 if [ $# -lt 2 ]; then
  echo "usage: $0 file dest [force]"
 else

  file=$1
  dest=$2
  force=${3:-noforce}

  sha=$(shasum -a 256 $file | cut -f 1 -d ' ')

 (curl -s https://baxx.dev/sha256/$BAXX_TOKEN/$sha -f >/dev/null 2>&1 \
    && [[ "$force" != "force" ]] \
    && echo SKIP $file .. already baxxed, use \"$0 $1 $2 force\" to force) || \
 curl -T $file https://baxx.dev/io/$BAXX_TOKEN/$dest

 fi
}


baxx_get() {
 if [ $# -ne 2 ]; then
  echo "usage: $0 file dest"
 else
  file=$1
  dest=$2
  curl https://baxx.dev/io/$BAXX_TOKEN/$file > $dest
fi
}


baxx_ls() {
 curl https://baxx.dev/ls/$BAXX_TOKEN/$*
}

---

check out https://github.com/jackdoe/baxx/tree/master/examples
for more examples

then simply do
% baxx_put example.txt /some/dir/example.txt
2918    Sun Mar 10 07:08:35 2019        /some/dir/example.txt@v2755

% baxx_get /some/dir/example.txt example.txt.dl

--
baxx.dev
`)

var EMAIL_VALIDATION = Parse(`Hi,

this is the verification link:

https://baxx.dev/verify/{{.ID}}

You can check the account status with:

curl -u {{.Email}} -XPOST https://baxx.dev/protected/status | json_pp

PS:
It is very likely that this email goes to the spam folder because it
is small and texty.. anyway, I hope it doesnt.

--
baxx.dev
`)

var EMAIL_PAYMENT_THANKS = Parse(`Hi,

Thanks for subscribing!
Even though the service is just in alpha state, it is much
appreciated!

If you want to cancel you have to do that in your paypal account,
or go to:
https://baxx.dev/unsub/{{ .PaymentID }}
which will redirect you there.

You can check the account status with:

curl -u {{.Email}} -XPOST https://baxx.dev/protected/status | json_pp

--
baxx.dev
`)

var EMAIL_WAIT_PAYPAL = Parse(`Hi,

Thanks for subscribing, it usually takes 1-2 minutes to receive
the notification from paypal, and then your account should be
enabled, if not please send me an email to jack@baxx.dev.

curl -u {{.Email}} -XPOST https://baxx.dev/protected/status | json_pp

--
baxx.dev
`)
