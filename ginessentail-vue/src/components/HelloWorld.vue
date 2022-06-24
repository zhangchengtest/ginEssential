<template>
  <div class="hello">
    <h1>{{ msg }}</h1>

 <div class="warp">
         <div class="myinput" >
            <input v-model="mymodel.tableName" />
        </div>
        <div>
        </div>
        <div class="mybutton" @click="gochange()">
            <button> 提交</button>
        </div>
      </div>

      <div class="warp">
         <div>
          <textarea class="mytextarea"  v-model="mymodel.originText">
    </textarea>
        </div>
        <div>
           >>
        </div>
        <div>
               <textarea class="mytextarea" v-model="responseText">
    </textarea>
        </div>
      </div>

  </div>
</template>

<script>
export default {
  name: 'HelloWorld',
  props: {
    msg: String,
  },
  data() {
    return {
      mymodel: {
        tableName: '',
        originText: '',
      },
      responseText: '',
    }
  },
  methods: {
    gochange() {
      // 请求api
      const api = 'http://127.0.0.1:8080/api/javatosql'
      this.axios.post(api, { ...this.mymodel }).then((res) => {
        const { data } = res.data
        console.log(data)
        this.responseText = data
        // 保存token

        // 跳转主页
      }).catch((err) => {
        console.log('err:', err.response)
        this.$bvToast.toast(err.response.data.msg, {
          title: '出错啦',
          variant: 'danger',
          solid: true,
        })
      })
      console.log('register')
    },
  },
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
