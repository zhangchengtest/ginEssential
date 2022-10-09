<template>
  <div class="hello">
    <h1>{{ msg }}</h1>

    <div class="warp">
      <div class="myinput">
        <input ref="file1" type="file" accept=".sql" @change="uploadFile">
        <input ref="file2" type="file" accept=".sql" @change="uploadFile2">
      </div>
      <div />
      <div class="mybutton" @click="submitAddFile()">
        <button> 提交</button>
      </div>
    </div>

    <div id="mycontent" class="warp" />
  </div>
</template>

<script>
import {
  compareFile
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
      myfile1: null,
      myfile2: null
    }
  },
  methods: {
    uploadFile() {
      this.myfile1 = this.$refs.file1.files[0]
    },
    uploadFile2() {
      this.myfile2 = this.$refs.file2.files[0]
    },
    submitAddFile() {
      var formData = new FormData()
      formData.append('file1', this.myfile1)
      formData.append('file2', this.myfile2)
         compareFile(formData).then((res) => {
        const { data } = res.data
        console.log(data)
        document.getElementById('mycontent').innerHTML = data
        // 保存token

        // 跳转主页
      }).catch((err) => {
        console.log('err:', err)
      })
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
