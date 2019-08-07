

function RefreshPage() {
    console.log("call RefreshPage")

    fetch("/get_records")
        .then(response => response.json())
        .then(function(response) {
            var errorTable = document.createElement("Table");
            errorTable.setAttribute("border", 2);
            //document.getElementsByTagName('body')[0].appendChild(errorTable);
            document.getElementById("table_section").appendChild(errorTable);
            var titleRow = document.createElement("tr");
            errorTable.appendChild(titleRow);

            var t1 = document.createElement("th");
            t1.innerHTML = "文本内容"
            var t2 = document.createElement("th");
            t2.innerHTML = "语言"
            var t3 = document.createElement("th");
            t3.innerHTML = "最新上报时间"
            var t4 = document.createElement("th");
            t4.innerHTML = "上报玩家"
            var t5 = document.createElement("th");
            t5.innerHTML = "文件名"
            var t6 = document.createElement("th");
            t6.innerHTML = "打包时间"

            titleRow.appendChild(t1);
            titleRow.appendChild(t2);
            titleRow.appendChild(t3);
            titleRow.appendChild(t4);
            titleRow.appendChild(t5);
            titleRow.appendChild(t6);
            
            for (var key in response) {
                //console.log(key, (response[key]))
                if (response[key]["Ignore"] != 1) {
                    CreateRow(errorTable, response[key]["Text_content"], response[key]["Language"], response[key]["Update_time"], response[key]["User_name"], response[key]["File_name"], response[key]["Version_str"])
                }  
            }
        })
}

function CreateRow(parentElement, textContent, language, updateTime, userName, fileName, versionStr) {
    var row = document.createElement("tr");
    var contentNode = document.createElement("td");
    contentNode.innerHTML = textContent;
    var langNode = document.createElement("td");
    langNode.innerHTML = language;
    var timeNode = document.createElement("td");

    var time = new Date(updateTime * 1000);
    timeNode.innerHTML = time.toString();

    var nameNode = document.createElement("td");
    nameNode.innerHTML = userName;

    var fileNameNode = document.createElement("td");
    fileNameNode.innerHTML = fileName;

    var verNode = document.createElement("td");
    verNode.innerHTML = versionStr; 

    row.appendChild(contentNode);
    row.appendChild(langNode);
    row.appendChild(timeNode);
    row.appendChild(nameNode);
    row.appendChild(fileNameNode);
    row.appendChild(verNode);

    parentElement.appendChild(row);
}