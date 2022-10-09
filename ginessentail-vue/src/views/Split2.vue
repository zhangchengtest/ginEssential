<template>
  <div class="hello">
    <h1>Split Image</h1>
    <form action="#">
      <p>
        <input id="row" type="hidden" name="row" value="3">

        <input id="column" type="hidden" name="column" value="3">

        <SingleImageUpload v-model="imageUrl" :cropper="true" />

        <!-- <el-upload
          ref="upload"
          class="pop-upload"
          action=""
          :file-list="fileList"
          :auto-upload="false"
          :multiple="true"
          :on-change="handleChange"
          :on-remove="handleRemove"
        >
          <el-button slot="trigger" size="small" type="primary">
            选取文件
          </el-button>
          <el-button style="margin-left: 10px;" size="small" type="success" @click="submitUpload">
            上传到服务器
          </el-button>
        </el-upload> -->
      </p>
    </form>

    <h2>Image Preview</h2>
    <div id="preview">
      try to drag an image here
    </div>

    <h2 style="margin-top: 60px;">
      Image Split Piece
    </h2>
    <div id="result" />
  </div>
</template>

<script>
import SingleImageUpload from '@/components/SingleImageUpload'
import {
  handleFile
} from '@/utils/split'
import {
  uploadSplitImages
} from '@/api/all'
export default {
  name: 'Split',
  components: {
    SingleImageUpload
  },
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
      imageUrl: '',
      mypieces: [],
      myfile: null
    }
  },

  methods: {
    handleChange(file) {
      console.log(file)
      this.myfile = file.raw
      handleFile(file.raw, this.mypieces)
    },
    // 上传服务器
    submitUpload() {
        console.log('Hello world')
        console.log(this.myfile)
        if (this.mypieces.length === 0) {
            return this.$message.warning('请选取文件后再上传')
        }
        const formData = new FormData()
        var i = 1
        this.mypieces.forEach((file) => {
            formData.append('piece' + i, file)
            i++
        })
        formData.append('file', this.myfile)
        uploadSplitImages(formData).then((res) => {
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
 #preview {
            width: 400px;
            height: 200px;
            border: 2px dashed #00f;
        }
        #preview img {
            width: 100%;
            height: 100%;
        }
        table {
            border-collapse: collapse;
        }
        td, th {
            padding: 0;
        }
        td img {
            display: block;
            padding: 2px;
            background: #fff;
        }

</style>
