import{_,d as w,m as h,h as C,u as v,i as b,j as E,k as T,s as V,l as k,r as p,o as $,c as S,a as s,t as y,b as u,w as l,n as A,E as m,p as D,e as B,f as F}from"./index.f1bf0a60.js";import{l as j}from"./left.f31b3d0a.js";const z=w({components:{selectLang:h},setup(){const e=C(),o=v();b();const t=E({name:"",password:"",loading:!1}),a=T("password"),g=()=>{a.value===""?a.value="password":a.value=""},f=()=>new Promise((n,i)=>{if(t.name===""){m.warning({message:"\u7528\u6237\u540D\u4E0D\u80FD\u4E3A\u7A7A",type:"warning"});return}if(t.password===""){m.warning({message:"\u5BC6\u7801\u4E0D\u80FD\u4E3A\u7A7A",type:"warning"});return}n(!0)});return{loginLeftPng:j,systemTitle:V,systemSubTitle:k,form:t,passwordType:a,passwordTypeChange:g,submit:()=>{f().then(()=>{t.loading=!0;let n={username:t.name,password:t.password};e.dispatch("user/register",n).then(async i=>{i.code===200&&(m.success({message:"\u6CE8\u518C\u6210\u529F",type:"success",showClose:!0,duration:2e3}),await o.push("/login"))})})}}}}),c=e=>(D("data-v-abcb0124"),e=e(),B(),e),I={class:"container"},N={class:"box"},L={class:"login-content-left"},P=["src"],R={class:"login-content-left-mask"},U={class:"box-inner"},M=c(()=>s("h1",null,"\u6B22\u8FCE\u6CE8\u518C",-1)),q=c(()=>s("i",{class:"sfont system-xingmingyonghumingnicheng"},null,-1)),G=c(()=>s("i",{class:"sfont system-mima"},null,-1)),H=F(" \u6CE8\u518C ");function J(e,o,t,a,g,f){const d=p("el-input"),n=p("el-button"),i=p("el-form");return $(),S("div",I,[s("div",N,[s("div",L,[s("img",{src:e.loginLeftPng},null,8,P),s("div",R,[s("div",null,y(e.$t(e.systemTitle)),1),s("div",null,y(e.$t(e.systemSubTitle)),1)])]),s("div",U,[M,u(i,{class:"form"},{default:l(()=>[u(d,{size:"large",modelValue:e.form.name,"onUpdate:modelValue":o[0]||(o[0]=r=>e.form.name=r),placeholder:e.$t("message.system.userName"),type:"text",maxlength:"50"},{prepend:l(()=>[q]),_:1},8,["modelValue","placeholder"]),u(d,{size:"large",ref:"password",modelValue:e.form.password,"onUpdate:modelValue":o[2]||(o[2]=r=>e.form.password=r),type:e.passwordType,placeholder:e.$t("message.system.password"),name:"password",maxlength:"50"},{prepend:l(()=>[G]),append:l(()=>[s("i",{class:A(["sfont password-icon",e.passwordType?"system-yanjing-guan":"system-yanjing"]),onClick:o[1]||(o[1]=(...r)=>e.passwordTypeChange&&e.passwordTypeChange(...r))},null,2)]),_:1},8,["modelValue","type","placeholder"]),u(n,{type:"primary",onClick:e.submit,style:{width:"100%"},size:"medium"},{default:l(()=>[H]),_:1},8,["onClick"])]),_:1})])])])}var Q=_(z,[["render",J],["__scopeId","data-v-abcb0124"]]);export{Q as default};