<ul class="nav nav-tabs" id="myTab">
  <li class="active"><a href="#home">主机列表</a></li>
</ul>

<div class="tab-content">
  <div role="tabpanel" class="tab-pane active" id="home" style="margin: 30px 0">
	<table border="0">
		<form  role="form" action="/ssh/index" method="post">
  		<tr>
    		<th width="700px">
      			<textarea class="form-control" rows="3" name="iplist" id="iplist" placeholder="IP:端口:账号:密码(英文字符)" style="width:700px;height:400px;">{{.iplist}}</textarea></td>
    		<td rowspan="2" style="padding:0 30px;align:left;vertical-align: top;">
				<textarea class="form-control" style="width:700px;height:430px">{{.rcmd}}</textarea>
			</td>
		</tr>

        <tr>
	        <th>
                <input type="text" style="width:700px" class="form-control" name="cmd" id="cmd" value="" placeholder="">
     			 <button type="submit" class="btn btn-primary" style="margin:15px 30px 15px 0">执行</button>
     			 {{.tips}}</br>
             </th>
		</tr>

        </form>
    </table>
  </div>
  <div role="tabpanel" class="tab-pane" id="profile">...</div>
</div>
