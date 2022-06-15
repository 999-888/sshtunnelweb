import{T as V}from"./index.3f34d0c0.js";import{d as C,L as A,k as y,H as L,_ as T,r as n,o as m,x as _,w as t,b as o,c as $,z as E,f as b,t as D,F as k,j as U,A as N,a as S,B as G,y as w,I as J,J as O,E as j}from"./index.79adc836.js";import{g as q}from"./remote.5f97b09a.js";import{f as z}from"./date.df283880.js";const B=[{value:!0,label:"\u662F"},{value:!1,label:"\u5426"}],I=C({components:{Layer:A},props:{layer:{type:Object,default:()=>({show:!1,title:"",showButton:!0})}},setup(e,l){const c=y(null),g=y(null);let f=y({id:"",username:"",isadmin:!1});const h={username:[{required:!0,message:"\u7528\u6237\u540D",trigger:"blur"}],isadmin:[{required:!0,message:"\u662F\u5426\u662Fadmin",trigger:"blur"}]};a();function a(){if(e.layer.row){f.value=JSON.parse(JSON.stringify(e.layer.row));const s=B.find(p=>p.label===e.layer.row.isadmin);s?f.value.isadmin=s.value:f.value.isadmin}}const r=y([]);function u(){q().then(s=>{r.value=s.data})}return u(),{radioData:B,connSelectData:r,form:f,rules:h,layerDom:g,ruleForm:c}},methods:{submit(){this.ruleForm&&this.ruleForm.validate(e=>{if(e){let l=this.form;this.layer.row&&this.updateForm(l)}else return!1})},updateForm(e){console.log(e),L(e).then(l=>{this.$message({type:"success",message:"\u7F16\u8F91\u6210\u529F"}),this.$emit("getTableData"),this.layerDom&&this.layerDom.close()})}}});function H(e,l,c,g,f,h){const a=n("el-input"),r=n("el-form-item"),u=n("el-radio"),s=n("el-radio-group"),p=n("el-form"),v=n("Layer",!0);return m(),_(v,{layer:e.layer,onConfirm:e.submit,ref:"layerDom"},{default:t(()=>[o(p,{model:e.form,rules:e.rules,ref:"ruleForm","label-width":"130px",style:{"margin-right":"30px"}},{default:t(()=>[o(r,{label:"\u7528\u6237\u540D: ",prop:"username"},{default:t(()=>[o(a,{modelValue:e.form.username,"onUpdate:modelValue":l[0]||(l[0]=d=>e.form.username=d),placeholder:"\u7528\u6237\u540D",disabled:""},null,8,["modelValue"])]),_:1}),o(r,{label:"\u662F\u5426\u662Fadmin: ",prop:"isadmin"},{default:t(()=>[o(s,{modelValue:e.form.isadmin,"onUpdate:modelValue":l[1]||(l[1]=d=>e.form.isadmin=d)},{default:t(()=>[(m(!0),$(k,null,E(e.radioData,d=>(m(),_(u,{key:d.value,label:d.value},{default:t(()=>[b(D(d.label),1)]),_:2},1032,["label"]))),128))]),_:1},8,["modelValue"])]),_:1})]),_:1},8,["model","rules"])]),_:1},8,["layer","onConfirm"])}var M=T(I,[["render",H]]);const P=C({name:"sshinfo",components:{Table:V,Layer:M},setup(){const e=y(!0),l=y([]),c=U({show:!1,title:"\u65B0\u589E",showButton:!0}),g=()=>{e.value=!0,J().then(a=>{let r=a.data;Array.isArray(r)&&r.forEach(u=>{const s=B.find(p=>p.value===u.isadmin);s?u.isadmin=s.label:u.isadmin}),l.value=r}).catch(()=>{l.value=[]}).finally(()=>{e.value=!1})},f=a=>{O({id:a}).then(r=>{j({type:"success",message:"\u5220\u9664\u6210\u529F"})})},h=a=>{c.title="\u7F16\u8F91\u6570\u636E",c.row=a,c.show=!0};return g(),{tableData:l,loading:e,layer:c,handleDel:f,handleEdit:h,getTableData:g,formatDate:z}}}),R={class:"layout-container"},K={class:"layout-container-table"},Q=b("\u66F4\u65B0");function W(e,l,c,g,f,h){const a=n("el-table-column"),r=n("el-tag"),u=n("el-button"),s=n("el-popconfirm"),p=n("Table"),v=n("Layer"),d=N("loading");return m(),$("div",R,[S("div",K,[G((m(),_(p,{ref:"table",data:e.tableData,onGetTableData:e.getTableData},{default:t(()=>[o(a,{prop:"username",label:"\u7528\u6237\u540D",align:"center"}),o(a,{prop:"ip",label:"\u7528\u6237\u8BBF\u95EEIP",align:"center"}),o(a,{prop:"isadmin",label:"\u662F\u5426\u662Fadmin",align:"center"}),o(a,{prop:"conn",label:"\u6388\u6743\u8BBF\u95EE\u7684\u670D\u52A1",align:"center"},{default:t(i=>[(m(!0),$(k,null,E(i.row.conn,F=>(m(),_(r,{key:F.id,size:"small",type:"success",style:{margin:"1px"}},{default:t(()=>[b(D(F.svcname),1)]),_:2},1024))),128))]),_:1}),o(a,{prop:"CreatedAt",label:"\u521B\u5EFA\u65F6\u95F4",align:"center"},{default:t(i=>[b(D(e.formatDate(i.row.CreatedAt)),1)]),_:1}),o(a,{prop:"UpdatedAt",label:"\u66F4\u65B0\u65F6\u95F4",align:"center"},{default:t(i=>[b(D(e.formatDate(i.row.UpdatedAt)),1)]),_:1}),o(a,{label:"\u64CD\u4F5C",align:"center",fixed:"right",width:"200"},{default:t(i=>[i.row.username!=="root"?(m(),_(u,{key:0,onClick:F=>e.handleEdit(i.row)},{default:t(()=>[Q]),_:2},1032,["onClick"])):w("",!0),o(s,{title:e.$t("message.common.delTip"),onConfirm:F=>e.handleDel(i.row.id)},{reference:t(()=>[i.row.username!=="root"?(m(),_(u,{key:0,type:"danger"},{default:t(()=>[b(D(e.$t("message.common.del")),1)]),_:1})):w("",!0)]),_:2},1032,["title","onConfirm"])]),_:1})]),_:1},8,["data","onGetTableData"])),[[d,e.loading]]),e.layer.show?(m(),_(v,{key:0,layer:e.layer,onGetTableData:e.getTableData},null,8,["layer","onGetTableData"])):w("",!0)])])}var ee=T(P,[["render",W]]);export{ee as default};
