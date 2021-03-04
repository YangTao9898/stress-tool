/**
 * 对外提供的方法
 */
$.InputDataComponent.addInputData = addInputData
$.InputDataComponent.getInputDataElement = getInputDataElement
$.InputDataComponent.inputDataIntChose = inputDataIntChose
$.InputDataComponent.inputDataFloatChose = inputDataFloatChose
$.InputDataComponent.inputDataStringChose = inputDataStringChose
$.InputDataComponent.inputDataNumEndClick = inputDataNumEndClick
$.InputDataComponent.getInputData = getInputData
$.InputDataComponent.setInputData = setInputData



/**
 * 生成一个输入数据的组件到 id 为 parentId 的节点下，内部存储数据的节点的 class 皆有一个前缀 classPrefix 用以区分各组件内部元素
 * @param parentId
 * @param classPrefix
 */
function addInputData(parentId, classPrefix) {
    var inputDataHTML = getInputDataElement(classPrefix)
    //$("#input-data-parent").html($("#input-data-parent").html() + inputDataHTML)
    $("#" + parentId).append($(inputDataHTML))
}

function getInputDataElement(itemClassPrefix) {
    var uniqueId = $.getUniqueId()
    return `
        <div class="` + itemClassPrefix + `input-data-top">
            <input type="hidden" class="` + itemClassPrefix + `input-data-type">
            <div style="display: inline-block; max-width: 30%">
                <div class="input-group input-group-lg top-pad-1">
                    ` + "<!-- 类型按钮 -->" + `
                    <div class="btn-group" style="float: left;">
                        <button type="button" class="btn btn-default dropdown-toggle btn-lg" data-toggle="dropdown"
                            aria-haspopup="true" aria-expanded="false"  style="width: 160px;">
                            <span class="dropdown-toggle-text">定义数据类型</span> <span class="caret"></span>
                        </button>
                        <ul class="dropdown-menu">
                            <li><a onclick="inputDataIntChose(this, '` + itemClassPrefix + `')">Int数据</a></li>
                            <li><a onclick="inputDataFloatChose(this, '` + itemClassPrefix + `')">Float数据</a></li>
                            <li><a onclick="inputDataStringChose(this, '` + itemClassPrefix + `')">字符串数据</a></li>
                            <li role="separator" class="divider"></li>
                            <li><a >取消</a></li>
                        </ul>
                    </div>
                    ` + "<!-- 类型长度 -->" + `
                    <span class="input-group-addon left-margin-1" style="float:left; width: 150px">长度（byte）<span class="must-fill">*</span></span>
                    <input type="text" class="form-control ` + itemClassPrefix + `input-data-length" style="width: 150px; float: left" name="Data" aria-describedby="sizing-addon1">
                </div>
            </div>
            <div style="display: inline-block; max-width: 50%">
                <div class="input-group input-group-lg top-pad-1">
                    ` + "<!-- 数据 -->" + `
                    <span class="input-group-addon left-margin-1">数据<span class="must-fill">*</span></span>
                    <input type="text" class="form-control ` + itemClassPrefix + `input-data-data" name="Data" aria-describedby="sizing-addon1">
                </div>
            </div>
            ` + "<!-- 大小端 -->" + `
            <div style="display: inline-block;" class="left-margin-2 ` + itemClassPrefix + `input-data-num-end-parent">
                <div class="input-group input-group-lg top-pad-1" style="height: 46px; font-size: 18px;">
                    <input type="checkbox" class="pointer ` + itemClassPrefix + `input-data-num-end" id="inputDataNumEnd` + uniqueId + `"
                        onclick="inputDataNumEndClick(this, '` + itemClassPrefix + `')" checked>
                    <label class="pointer ` + itemClassPrefix + `input-data-num-end-label" style="margin-left: 5px; font-weight: normal;" for="inputDataNumEnd` + uniqueId + `">大端</label>
                </div>
            </div>
            <div style="display: inline-block; max-width: 10%" class="left-margin-2">
                <div class="input-group input-group-lg top-pad-1">
                    ` + "<!-- 数据 -->" + `
                    <button style="display: inline-block;"
                        class="btn btn-default btn-lg" onclick="deleteElement(this)">删除</button>
                </div>
            </div>
        </div>`
}

function deleteElement(obj) {
    $(obj.parentNode.parentNode.parentNode).remove()
}

function inputDataIntChose(obj, classPrefix) {
    var cprefix = "." + classPrefix
    var p = $(obj).parents(cprefix + "input-data-top")[0]
    var inputLength = $(p).find(cprefix + "input-data-length")[0]
    $(inputLength).attr("placeholder", "只能为1,2,4,8")
    $(inputLength).val("")
    $(inputLength).removeAttr("disabled")
    var inputType = $(p).children(cprefix + "input-data-type")[0]
    $(inputType).val("0")
    $(p).find(".dropdown-toggle-text")[0].innerText = "Int数据"
    // 大小端复选框
    $(p).find(cprefix + "input-data-num-end-parent").removeClass("hidden")

    var inputData = $(p).find(cprefix + "input-data-data")[0]
    //$(inputData).unbind('input', computeBytes)
    $(inputData).unbind('blur')
}

function inputDataFloatChose(obj, classPrefix) {
    var cprefix = "." + classPrefix
    var p = $(obj).parents(cprefix + "input-data-top")[0]
    var inputLength = $(p).find(cprefix + "input-data-length")[0]
    $(inputLength).attr("placeholder", "只能为4,8")
    $(inputLength).val("")
    $(inputLength).removeAttr("disabled")
    var inputType = $(p).children(cprefix + "input-data-type")[0]
    $(inputType).val("1")
    $(p).find(".dropdown-toggle-text")[0].innerText = "Float数据"
    // 大小端复选框
    $(p).find(cprefix + "input-data-num-end-parent").removeClass("hidden")

    var inputData = $(p).find(cprefix + "input-data-data")[0]
    //$(inputData).unbind('input', computeBytes)
    $(inputData).unbind('blur')
}

/*function computeBytes(event) {
    var objlength = event.data.inputLength
    var objdata = event.data.inputData
    var data = $(objdata).val()

    for (var i = 0; i < data.length; i++) {
        console.log(data.charCodeAt(i))
    }
    //console.log(atob(data))
    $(objlength).val()
}*/
function inputDataStringChose(obj, classPrefix) {
    var cprefix = "." + classPrefix
    var p = $(obj).parents(cprefix + "input-data-top")[0]
    var inputLength = $(p).find(cprefix + "input-data-length")[0]
    $(inputLength).attr("placeholder", "将会自动补充")
    $(inputLength).val("")
    $(inputLength).attr("disabled", "")
    var inputType = $(p).children(cprefix + "input-data-type")[0]
    $(inputType).val("2")
    $(p).find(".dropdown-toggle-text")[0].innerText = "字符串数据"
    // 大小端复选框
    $(p).find(cprefix + "input-data-num-end-parent").addClass("hidden")

    var inputData = $(p).find(cprefix + "input-data-data")[0]
    var lastValue = inputData.value
    $(inputData).bind('blur', function () {
        // input值未变时，不发送请求
        if (lastValue != inputData.value || $(inputLength).val() == "") {
            lastValue = inputData.value
            $.ajax({
                url: "/GetStrBytes",
                type: "POST",
                data: inputData.value,
                success: function(data) {
                    $(inputLength).val(data)
                },
                error: function (err) {
                    $.showRequestErr(err)
                }
            })
        }
    })
}

function inputDataNumEndClick(obj, classPrefix) {
    var cprefix = "." + classPrefix
    var parent = $(obj).parent()
    var checkbox = $(parent).find(cprefix + "input-data-num-end")
    $(parent).find(cprefix + "input-data-num-end-label").text($(checkbox).is(":checked") ? "大端" : "小端")
}

/**
 * 返回输入的数据
 * @param inputDataParentId 父节点id
 * @param classPrefix class 名称前缀
 * @returns {[]}
 */
function getInputData(inputDataParentId, classPrefix, showErrBox) {
    var cPrefix = "." + classPrefix
    var inputDataMap = []
    var inputDataParent = $("#" + inputDataParentId)
    var inputDataArr = $(inputDataParent).children(cPrefix + "input-data-top")
    var ret = {}
    ret.inputDataMap = null
    ret.hasErr = true
    for (var i = 0; i < inputDataArr.length; i++) {
        var node = inputDataArr[i]
        var type = $(node).children(cPrefix + "input-data-type")[0].value
        var inputLength = $(node).find(cPrefix + "input-data-length")[0].value
        var inputData = $(node).find(cPrefix + "input-data-data")[0].value
        var isBigEnd = $(node).find(cPrefix + "input-data-num-end")[0].checked
        var obj = {}
        if (type == "") {
            ret.errMsg = "数据类型不能为空"
            if (showErrBox == true) {
                $.showAlertBox("请检查", "数据类型不能为空", true)
            }
            return ret
        }
        if (inputLength == "" && type != "2") {
            ret.errMsg = "数据长度不能为空"
            if (showErrBox == true) {
                $.showAlertBox("请检查", "数据长度不能为空", true)
            }
            return ret
        }
        if (inputData == "") {
            ret.errMsg = "数据不能为空"
            if (showErrBox == true) {
                $.showAlertBox("请检查", "数据不能为空", true)
            }
            return ret
        }

        obj.type = type
        obj.length = inputLength
        obj.data = inputData
        obj.isBigEnd = isBigEnd
        inputDataMap.push(obj)
    }
    ret.hasErr = false
    ret.inputDataMap = inputDataMap
    return ret
}

/**
 * 给输入组件设置值
 * @param data [{"type":"", "length":"", "data":"", "isBigEnd":""}]
 */
function setInputData(parentId, classPrefix, data) {
    var parentNode = $("#" + parentId)
    parentNode.html("")

    var cPrefix = "." + classPrefix
    for (var i = 0, length = data.length; i < length; i++) {
        var item = data[i]
        parentNode.append($(getInputDataElement(classPrefix)))
        var node = $(parentNode).children(cPrefix + "input-data-top")[i] // 获取到刚才新增的最后一个节点

        var type = $(node).children(cPrefix + "input-data-type")[0]
        var inputLength = $(node).find(cPrefix + "input-data-length")[0]
        var inputData = $(node).find(cPrefix + "input-data-data")[0]
        var isBigEnd = $(node).find(cPrefix + "input-data-num-end")[0]

        switch(item.type) {
            case "0":
                $.InputDataComponent.inputDataIntChose(type, classPrefix)
                break;
            case "1":
                $.InputDataComponent.inputDataFloatChose(type, classPrefix)
                break;
            case "2":
                $.InputDataComponent.inputDataStringChose(type, classPrefix)
                break;
            default:
                $.showAlertBox("错误", "未知数据类型", true);
                return
        }
        inputLength.value = item.length
        inputData.value = item.data
        $(isBigEnd).prop("checked", item.isBigEnd)
        $.InputDataComponent.inputDataNumEndClick(type, classPrefix)
    }
}