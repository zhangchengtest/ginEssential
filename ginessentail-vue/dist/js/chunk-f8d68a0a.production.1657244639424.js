(window["webpackJsonp"]=window["webpackJsonp"]||[]).push([["chunk-f8d68a0a"],{"165d":function(e,t,o){"use strict";o.r(t);var n=function(){var e=this,t=e.$createElement,o=e._self._c||t;return o("div",{staticClass:"hello"},[o("h1",[e._v(e._s(e.msg))]),o("div",{staticClass:"warp"},[o("div",{staticClass:"myinput"},[e._v(" "+e._s(e.mymodel.chapter)+" - "+e._s(e.mymodel.readCount)+" ")]),o("div"),o("div",{staticClass:"mybutton",on:{click:function(t){return e.gochange()}}},[o("button",[e._v(" 提交")])])]),o("div",{staticClass:"warp"},[o("div",[o("textarea",{directives:[{name:"model",rawName:"v-model",value:e.mymodel.originText,expression:"mymodel.originText"}],staticClass:"mytextarea",domProps:{value:e.mymodel.originText},on:{input:function(t){t.target.composing||e.$set(e.mymodel,"originText",t.target.value)}}})]),o("div",[e._v(" >> ")]),o("div",[o("textarea",{directives:[{name:"model",rawName:"v-model",value:e.responseText,expression:"responseText"}],staticClass:"mytextarea",domProps:{value:e.responseText},on:{input:function(t){t.target.composing||(e.responseText=t.target.value)}}})])])])},a=[],i=o("8c63"),s={name:"HelloWorld",props:{msg:{type:String,default:""}},data:function(){return{mymodel:{chapter:"",originText:"",question:"",readCount:0},responseText:""}},mounted:function(){var e=this;Object(i["b"])().then((function(t){var o=t.data.data;console.log(o),e.mymodel.chapter=o.article.chapter,e.mymodel.readCount=o.article.readCount,e.mymodel.originText=o.article.content,e.mymodel.question=o.article.question})).catch((function(e){console.log("err:",e)}))},methods:{gochange:function(){alert(this.mymodel.question)}}},r=s,l=(o("adf8"),o("2877")),c=Object(l["a"])(r,n,a,!1,null,"439919d5",null);t["default"]=c.exports},adf8:function(e,t,o){"use strict";o("be80")},be80:function(e,t,o){}}]);