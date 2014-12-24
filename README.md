go-web-proxy
============

A simple web based proxy written in Go. Can be easily hosted on [OpenShift](https://openshift.redhat.com/app/login)

currently provides to options,

1] this is useful for a single file. It only pipes the response.

    ` <hostname>/p/?target=http://foo.com/bar.png`

2] this is useful to view a web page that contains other static images, css, js files

    `<hostname>/t/?target=http://foo.com/bar/`


how to host on OpenShift
------------------------

- fork this repo
- register an account on openshift.redhat.com, follow standard account activation procedure
- go to **Settings** tab and add public key for ssh to work
- go back to **Applications** tab, click **Add Application**
- scroll down and select Go language under Other types
- fill in the form, enter name/domain pair
- enter git repo address for forked repo. enter **master** as branch name
- keep default for rest but know that if you select **No scaling** then you won't be able to scale later
- Click **Create Application**
- Done


Add a bookmarklet to quickly use this proxy
-------------------------------------------

- Add a new bookmark in your browser
- name it anything, say **proxy**
- use this as its URL

    `javascript:location.href='http://proxy-gopherlang.rhcloud.com/t/?target='+location.href`

- now when you need to use this proxy, simply click the bookmark and it should work


License
-------

MIT license


