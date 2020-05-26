
function classQuery() {
    classinfo = document.getElementById("classinfo").value;
    var text;
    var xmlhttp;
    if (window.XMLHttpRequest){
        //  IE7+, Firefox, Chrome, Opera, Safari 浏览器执行代码
        xmlhttp=new XMLHttpRequest();
    }else{
        // IE6, IE5 浏览器执行代码
        xmlhttp=new ActiveXObject("Microsoft.XMLHTTP");
    }
    xmlhttp.onreadystatechange=function()
    {
        if (xmlhttp.readyState==4 && xmlhttp.status==200)
        {
            document.getElementById("classTemp").innerHTML=xmlhttp.responseText;
        }
    }
    xmlhttp.open("post","/template",true)
    xmlhttp.setRequestHeader("Content-type","application/x-www-form-urlencoded");
    text = "input="+classinfo;

    text = text+"&";
    text = text + "type=class";

    xmlhttp.send(text);
}


function goodsQuery() {
    goodsinfo = document.getElementById("goodsInfo").value;
    var xmlhttp;
    if (window.XMLHttpRequest){
        //  IE7+, Firefox, Chrome, Opera, Safari 浏览器执行代码
        xmlhttp=new XMLHttpRequest();
    }else{
        // IE6, IE5 浏览器执行代码
        xmlhttp=new ActiveXObject("Microsoft.XMLHTTP");
    }
    xmlhttp.onreadystatechange=function()
    {
        if (xmlhttp.readyState==4 && xmlhttp.status==200)
        {
            document.getElementById("goodsTemp").innerHTML=xmlhttp.responseText;
            //$("goodsTemp").html(xmlhttp.responseText);
        }
    }
    xmlhttp.open("post","/template",true)
    xmlhttp.setRequestHeader("Content-type","application/x-www-form-urlencoded");
    text = "input="+goodsinfo;
    text += "&";
    text += "type=goods";
    xmlhttp.send(text);
}

