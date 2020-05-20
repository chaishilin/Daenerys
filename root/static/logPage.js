getNewCapId()
var capId;

//log function
function getNewCapId() {
    $.get("/captcha/newId",
        function (data) {
            //console.log("change ",capId,"to ",data)
            capId = data;
            document.getElementById("cap-image").src = "/captcha/"+capId+".png"
        }
    )
}
function reloadVerify() {
    var cap,src;
    var reload = "reload=" + (new Date()).getTime();
    cap = document.getElementById("cap-image");
    src=cap.src;
    var p = src.indexOf('?');
    if (p >= 0) {
        src = src.substr(0, p);
    }
    cap.src = src + "?" + reload
}
function capVerify() {
    var vcode,text;
    vcode = document.getElementById("log-vcode").value;
    $.post("/process",
        {
            captchaId:capId,
            captchaSolution:vcode
        },function (data) {
            console.log(data)
            if(data!="true"){
                getNewCapId()
                text = "验证码错误";
            }else {
                text = "登录成功";
                location="/hello"
            }
            document.getElementById("logAlert").innerHTML = text;
        })
}
function logVerify() {
    //console.log("pwdVerify")
    var pwd,name,text;
    pwd = document.getElementById("log-pwd").value;
    name = document.getElementById("log-username").value;
    $(document).ready(function () {
        $.post("/",
            {
                name:name,
                pwd:pwd
            },
            function (data) {
                //console.log(data)
                if (data == "true"){
                    capVerify()
                }else{
                    text = "请输入正确的用户名和密码";
                }
                document.getElementById("logAlert").innerHTML = text;
            })
    })
}



var pwdReady = false;
var emailReady = false;
function regist() {
    var username,passwd,passwd2,email,vcode;
    username = document.getElementById("regist-username").value;
    passwd = document.getElementById("regist-pwd").value;
    email = document.getElementById("regist-email").value;
    if (pwdReady && emailReady){
        usernameVerify(username,passwd,email);
    }else{
        console.log("no")
    }
}
function usernameVerify(username,passwd,email) {
    $.post('/regist',{
            username:username,
            passwd:passwd,
            email:email
        },function (data) {
            if (data == "exist") {
                text = "用户名已被占用";
            } else {
                text = "注册成功";
            }
            document.getElementById("registAlert").innerText = text;
        }
    )
}
function regexEmail(email) {
    re = /^[0-9a-zA-Z]+([\.-_0-9a-zA-Z])*@([a-zA-Z0-9]+\.)+[a-zA-Z0-9]{2,4}$/
    if (re.test(email)){
        document.getElementById("emailHint").innerText = "";
        emailReady = true;
    }else {
        document.getElementById("emailHint").innerText = "邮箱地址无效！";
    }
}
function pwdConfirm() {
    passwd = document.getElementById("regist-pwd").value;
    passwd2 = document.getElementById("regist-pwd2").value;
    if (passwd != passwd2){
        document.getElementById("pwdHint").innerHTML="密码不一致！";

    }else{
        document.getElementById("pwdHint").innerHTML="";
        pwdReady = true;
    }
}
