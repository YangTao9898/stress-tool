// taskDetail 导出专用变量
var TaskDetal = {}
// inputDataComponent 导出专用变量
$.InputDataComponent = {}

var alertBoxCallBackFunc = null
$.alertBoxCallback = function () {
    if (alertBoxCallBackFunc != null) {
        alertBoxCallBackFunc()
    }
}

/**
 * 弹出框，需要引入 globalAlertBox.html 组件且同一页面只允许出现一个该组件，确定按钮没有回调函数 okBtnCallBack 则不会显示确定按钮
 * @param title string 弹出框的标题
 * @param content string 弹出框的内容
 * @param isErrAlert bool 如果为 true，则展示错误弹出框的样式
 * @param okBtnCallBack function 确认按钮按了之后的回调函数
 * @param okBtnMsg string 确认按钮的按键提示
 */
$.showAlertBox = function (title, content, isErrAlert, okBtnCallBack, okBtnMsg) {
    var headerStyle = ' style="background-color: #C4EDAF"'
    if (isErrAlert == true) {
        headerStyle = ' style="background-color: #FCA2A2"'
    }
    if (title == "") {
        title = "&nbsp;"
    }
    if (title == null) {
        if (isErrAlert == true) {
            title = "一个错误发生了"
        } else {
            title = "&nbsp;"
        }
    }
    if (okBtnMsg == "" || okBtnMsg == null) {
        okBtnMsg = "确认"
    }
    var showOkBtn = okBtnCallBack != null
    var node =
        '    <div id="global-alertbox-parent-div">' +
        '       <!-- Button trigger modal -->' +
        '       <button type="button" id="global-alertbox-trigger-btn" class="btn btn-primary hidden" data-toggle="modal" data-target="#globalAlertBox">' +
        '           Launch demo modal' +
        '       </button>' +
        '       <!-- Modal -->' +
        '       <div class="modal fade" id="globalAlertBox" tabindex="-1" role="dialog" aria-labelledby="globalAlertBoxLabel" aria-hidden="true">' +
        '           <div class="modal-dialog">' +
        '               <div class="modal-content">' +
        '                   <div class="modal-header"' + headerStyle + '>' +
        '                       <h5 class="modal-title" id="globalAlertBoxLabel">' + title + '</h5>' +
        '                   </div>' +
        '                   <div class="modal-body">' + content +
        '                   </div>' +
        '                   <div class="modal-footer">' +
        '                       <button type="button" class="btn btn-secondary" data-dismiss="modal">关闭</button>' +
        (showOkBtn ? ('           <button type="button" class="btn btn-primary" data-dismiss="modal" onclick="$.alertBoxCallback()">' + okBtnMsg + '</button>') : '') +
        '                   </div>' +
        '               </div>' +
        '           </div>' +
        '       </div>' +
        '   </div>'
    // 删除原来的弹出框
    $("#global-alertbox-parent-div").remove()
    // 模态框出现后，会自动添加 .modal-backdrop.fade.in，故也同时删除
    $(".modal-backdrop.fade.in").remove()
    $(document.body).append($(node))
    alertBoxCallBackFunc = okBtnCallBack
    $("#global-alertbox-trigger-btn").click()
}

$.loadingStart = function () {
    var node = `
    <div id="global-load-div">
        <div class="modal-backdrop show" style="opacity:.5;">
        </div>
        <div style="position: fixed; top: 0; bottom: 0; left: 0; right: 0; z-index: 1040">
            <div style="position: absolute; width: 100%; height: 100%;">
                <div class="loader-05" style="position: relative; top: 50%; left: 50%; margin: -1em 0 0 -1em;"></div>
            </div>
        </div>
    </div>
    `
    $("#global-load-div").remove()
    $(document.body).append($(node))
}

$.loadingEnd = function () {
    $("#global-load-div").remove()
}

var maskClickCallbackFunc = null
$.maskClickCallback = function () {
    if (maskClickCallbackFunc != null) {
        maskClickCallbackFunc()
    }
}
/**
 * 弹出一个半透明的遮罩
 * @param parent 给弹出的遮罩指定父级对象
 * @param z_index 遮罩的 z-index 层级
 * @param maskClickCallback 点击遮罩时触发的回调函数
 */
$.maskShow = function (parent, z_index, maskClickCallback) {
    var p = document.body
    var zIndex = 1
    if (z_index != "" && z_index != null) {
        zIndex = z_index
    }
    if (parent != null) {
        p = parent
    }
    var node = `
    <div id="global-mask-div" onclick="$.maskClickCallback()">
        <div class="modal-backdrop show" style="opacity:.5; z-index: ` + zIndex + `;">
        </div>
        </div>
    </div>
    `
    maskClickCallbackFunc = maskClickCallback
    $("#global-mask-div").remove()
    $(p).append($(node))
}

$.maskRemove = function () {
    $("#global-mask-div").remove()
}

/**
 * 打印 json 对象 data 是否是服务器正常的返回结果，如果不正常且为通用错误码则弹出错误框
 * @param data
 * @returns {boolean} false: 结果正常，true：结果不正常
 */
$.checkResultErr = function (data) {
    if (data.resultCode == "500") {
        $.showAlertBox("错误", "错误原因：" + data.errMsg, true, null, null)
    }
    return data == null || data.resultCode != "0000"
}

/**
 * 打印 ajax 访问错误的方法，传入 XMLHttpRequestErr，err 函数的第一个参数
 * @param XMLHttpRequestErr
 */
$.showRequestErr = function (XMLHttpRequestErr) {
    $.showAlertBox("连接服务器错误，请稍后访问", "错误原因：" + XMLHttpRequestErr.statusText, true, null, null)
}

/**
 * 格式化字符串
 * 用法，字符串对象直接打点调用，如
 //两种调用方式
 var template1="我是{0}，今年{1}了";
 var result1=template1.format("loogn",22);

 var template2="我是{name}，今年{age}了";
 var result2=template2.format({name:"loogn",age:22});

 //两个结果都是"我是loogn，今年22了"
 * @param args
 * @returns {String}
 */
String.prototype.format = function(args) {
    var result = this;
    if (arguments.length < 1) {
        return result;
    }

    var data = arguments; // 如果模板参数是数组
    if (arguments.length == 1 && typeof (args) == "object") {
        //如果模板参数是对象
        data = args;
    }
    for (var key in data) {
        var value = data[key];
        if (undefined != value) {
            result = result.replace("{" + key + "}", value);
        }
    }
    return result;
}

/**
 * 生成随机数
 * @param digit 随机数的位数
 * @returns {number}
 */
$.getRandomNum = function (digit) {
    var num = Math.random()
    retNum = num * Math.pow(10, digit)
    retNum = Math.floor(retNum)
    var len = (retNum + "").length
    if (len < digit) {
        retNum = num * Math.pow(10, digit + (digit - len))
        retNum = Math.floor(retNum)
    }
    return retNum
}

$.getUniqueId = function () {
    return new Date().getTime() + "-" + $.getRandomNum(6)
}