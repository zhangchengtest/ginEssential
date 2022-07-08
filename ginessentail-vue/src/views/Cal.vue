<template>
  <div class="hello">
    <h1>{{ msg }}</h1>

    <div class="warp">
      <div class="myinput">
        <el-upload
            class="pop-upload"
            ref="upload"
            action=""
            :file-list="fileList"
            :auto-upload="false"
            :multiple="true"
            :on-change="handleChange"
            :on-remove="handleRemove"
    >
     <el-button slot="trigger" size="small" type="primary">选取文件</el-button>
     <el-button style="margin-left: 10px;" size="small" type="success" @click="submitUpload">上传到服务器</el-button>

     <el-button style="margin-left: 10px;" size="small" type="primary" @click="recal"> 重新计算</el-button>
    </el-upload>

      </div>
      <div />
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
  readorc
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
        tableName: '',
        originText: ''
      },
      responseText: '',
      fileList: []
    }
  },
  methods: {
    handleChange(file, fileList) {
        this.fileList = fileList
    },
    handleRemove(file, fileList) {
        this.fileList = fileList
    },
    recal() {
        const myArray = this.mymodel.originText.split('\n')
        console.log(myArray)
        var sum = 0
        for (var i = 0; i < myArray.length; i++) {
            var ret = myArray[i].trim()
            if (myArray[i] && ret !== '') {
              sum += parseInt(myArray[i])
            }
        }
        this.responseText = sum
    },
    // 上传服务器
    submitUpload() {
        console.log('Hello world')
        if (this.fileList.length === 0) {
            return this.$message.warning('请选取文件后再上传')
        }
        const formData = new FormData()
        this.fileList.forEach((file) => {
            formData.append('file', file.raw)
        })
        readorc(formData).then((res) => {
          const { data } = res.data

          this.mymodel.originText = data
          const myArray = data.split('\r\n')
          var sum = 0
          for (var i = 0; i < myArray.length; i++) {
              if (myArray[i]) {
                sum += parseInt(myArray[i])
                console.log(myArray[i])
              }
          }
          this.responseText = sum
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
  width: 100%;
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
    margin-top: 10px;
}

</style>
