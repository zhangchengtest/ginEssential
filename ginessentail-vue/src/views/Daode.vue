<template>
  <div class="hello">
    <h1>{{ msg }}</h1>

    <div class="warp">
      <div class="myinput">
        {{mymodel.chapter}} - {{mymodel.readCount}}
      </div>
      <div />
      <div class="mybutton" @click="gochange()">
        <button> 提交</button>
      </div>
    </div>

    <div class="warp">
      <div>
        <textarea v-model="mymodel.originText" class="mytextarea" />
      </div>
      <div>
        >>
      </div>
      <div>
        <textarea v-model="responseText" class="mytextarea" />
      </div>
    </div>
  </div>
</template>

<script>
import {
  daode
} from '@/api/all'

export default {
  name: 'HelloWorld',
  props: {
    msg: {
      type: String,
      default: ''
    }
  },
  data() {
    return {
      mymodel: {
        chapter: '',
        originText: '',
        question: '',
        readCount: 0
      },
      responseText: ''
    }
  },
  mounted() {
     // 请求api
      daode().then((res) => {
        const { data } = res.data
        console.log(data)
        this.mymodel.chapter = data.article.chapter
        this.mymodel.readCount = data.article.readCount
        this.mymodel.originText = data.article.content
        this.mymodel.question = data.article.question
        // 保存token

        // 跳转主页
      }).catch((err) => {
        console.log('err:', err)
      })
  },
  methods: {
    gochange() {
      alert(this.mymodel.question)
    }
  }
}
</script>
<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped lang="scss">
h3 {
  margin: 40px 0 0;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}
.mytextarea{
  width: 500px;
  height: 600px;
}

.myinput{
   text-align: left;
   display: flex;
}

.mybutton{
  text-align: right;
  display: flex;
 justify-content: flex-end;
}

.warp {
    width: 100%;
    display: flex;
    flex-direction: row
}

.warp>div {
  // border: solid red 1px;
    flex: 1;
    text-align: center;
    margin-top: 10px;
}
</style>
