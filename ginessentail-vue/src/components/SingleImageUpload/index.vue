<template>
  <el-row class="row-box" type="flex" align="middle">
    <el-col :span="8">
      <el-button type="primary" size="mini" @click="upLoad">
        点击上传
      </el-button>
      <transition name="fade">
        <p v-if="!url" class="tip">
          {{ typeTip }}
        </p>
      </transition>
    </el-col>
    <el-col :span="16">
      <div v-if="url">
        <el-image
          v-if="type==='logo'|| type==='ico'"
          class="img-box"
          :src="url"
          :preview-src-list="[url]"
          :z-index="9999"
        />
        <div v-else>
          <el-row :gutter="20" type="flex" align="middle">
            <el-col :span="16">
              <span>{{ fileName }}</span>
            </el-col>
            <el-col :span="8">
              <el-button type="success" size="mini" @click="downFileLoad">
                下载
              </el-button>
            </el-col>
          </el-row>
        </div>
      </div>
    </el-col>
    <el-dialog
      title="提示"
      :visible.sync="dialogVisible"
      width="30%"
      append-to-body
    >
      <div class="cropper-box">
        <vueCropper
          ref="cropper"
          :img="option.img"
          :output-type="option.outputType"
          :info="option.info"
          :full="option.full"
          :can-move="option.canMove"
          :auto-crop-width="size.autoCropWidth"
          :auto-crop-height="size.autoCropHeight"
          :can-move-box="option.canMoveBox"
          :original="option.original"
          :auto-crop="option.autoCrop"
          :fixed="option.fixed"
          :fixed-number="option.fixedNumber"
          :center-box="option.centerBox"
          :info-true="option.infoTrue"
          :fixed-box="option.fixedBox"
        />
      </div>
      <span slot="footer" class="dialog-footer">
        <el-button @click="dialogVisible = false">取 消</el-button>
        <el-button type="primary" @click="submit">确 定</el-button>
      </span>
    </el-dialog>
  </el-row>
</template>

<script>
import OSS from 'ali-oss'
import { authTemp, tempFormal } from '@/api/common'
import { VueCropper } from 'vue-cropper'
export default {
  name: 'SingleImageUpload',
  components: {
    VueCropper
  },
  model: {
    prop: 'url',
    event: 'change'
  },
  props: {
    url: {
      type: String,
      required: false,
      default: ''
    },
    accept: {
      type: String,
      required: false,
      default: 'image/png, image/jpeg, image/jpg'
    },
    type: {
      type: String,
      required: false,
      default: 'logo'
    },
    cropper: {
      type: Boolean,
      default: false
    },
    size: {
      type: Object,
      default() {
        return {
          autoCropWidth: 200, // 默认生成截图框宽度
          autoCropHeight: 200 // 默认生成截图框高度
        }
      }
    }
  },
  data() {
    return {
      ossClient: null,
      uploadUrl: '',
      fileName: '',
      lastModified: '',
      dialogVisible: false,
      percentage: 0,
      option: {
        img: '', // 裁剪图片的地址
        info: true, // 裁剪框的大小信息
        outputSize: 1, // 裁剪生成图片的质量
        outputType: 'png', // 裁剪生成图片的格式
        canScale: true, // 图片是否允许滚轮缩放
        autoCrop: true, // 是否默认生成截图框
        autoCropWidth: 339, // 默认生成截图框宽度
        autoCropHeight: 151, // 默认生成截图框高度
        fixedBox: true, // 固定截图框大小 不允许改变
        fixed: false, // 是否开启截图框宽高固定比例
        full: false, // 是否输出原图比例的截图
        canMoveBox: false, // 截图框能否拖动
        original: true, // 上传图片按照原始比例渲染
        centerBox: true, // 截图框是否被限制在图片里面
        infoTrue: false // true 为展示真实输出图片宽高 false 展示看到的截图框宽高
      }
    }
  },
  computed: {
    typeTip() {
      switch (this.type) {
        case 'ico':
          return '请选择icon格式文件'
        case 'logo':
          return '请选择png,jpeg,jpg格式文件'
        case 'document':
          return '请选择pdf, doc, docx, rar, zip格式文件'
        default:
          return ''
      }
    }
  },
  mounted() {
    this.authTemp()
  },
  methods: {
    upLoad() {
      const fileDom = document.createElement('input')
      fileDom.type = 'file'
      fileDom.accept = this.accept
      fileDom.click()
      fileDom.onchange = (e) => {
        const file = e.path[0].files[0]
        if (file.size > 200 * 1024 * 1024) {
          fileDom.value = ''
          this.$message.error('文件最大不超过200m')
          return
        }
        this.lastModified = new Date().getTime()
        this.fileName = file.name
        const storeAs = this.uploadUrl + this.lastModified + file.name.substr(file.name.lastIndexOf('.'))
        if (this.cropper) {
          const reader = new FileReader()
          reader.readAsDataURL(e.path[0].files[0])
          reader.onload = (event) => {
            this.option.img = event.target.result
          }
          this.dialogVisible = true
        } else {
          this.multipartUpload(storeAs, file)
        }
      }
    },
    submit() {
      this.$refs.cropper.getCropBlob((data) => {
        const storeAs = this.uploadUrl + this.lastModified + this.fileName.substr(this.fileName.lastIndexOf('.'))
        this.multipartUpload(storeAs, data)
        this.dialogVisible = false
      })
    },
    multipartUpload(storeAs, file) {
      // 参数： storeAs为文件名称，file为文件
      const loading = this.$loading({
        lock: true,
        text: '上传中...',
        spinner: 'el-icon-loading',
        background: 'rgba(0, 0, 0, 0.7)'
      })
      this.percentage = 0
      this.ossClient.multipartUpload(storeAs, file, {
        progress: (p) => {
          this.percentage = Math.trunc(p * 100)
        }
      }).then((results) => {
        this.tempFormal(results.name, this.type)
      }).catch((err) => {
        console.log(err)
      }).finally(() => {
        this.$nextTick(() => { // 以服务的方式调用的 Loading 需要异步关闭
          loading.close()
        })
      })
    },
    tempFormal(relativePath, serveCode) {
      const data = {
        relativePath,
        serveCode
      }
      console.log(data)
      tempFormal(data).then(res => {
        this.$emit('change', res.data.absolutePath)
        if (this.fileName) {
          this.$emit('getFileName', this.fileName)
        }
      }).catch((err) => {
        console.log(err)
      })
    },
    authTemp() {
      authTemp().then(res => {
        this.ossClient = new OSS({
          accessKeyId: res.data.data.accessKeyId,
          accessKeySecret: res.data.data.accessKeySecret,
          stsToken: res.data.data.token,
          bucket: res.data.data.bucketName
        })
        this.uploadUrl = res.data.data.uploadUrl
      }).catch((err) => {
        console.log(err)
      })
    },
    downFileLoad() {
      window.location.href = this.url
    }
  }
}
</script>

<style scoped>
.img-box{
  width: 100px;
  height: 100px;
  border: 1px dashed #d9d9d9;
  border-radius: 5px;
}
.img-box:hover{
  border-color: #409eff;
}
.row-box{
  position: relative;
}
.tip{
  position: absolute;
  bottom: -30px;
  color: #F56C6C;
  font-size: 12px;
}
.fade-enter-active, .fade-leave-active {
  transition: all .3s;
}
.fade-enter, .fade-leave-to {
  opacity: 0;
  transform: translateY(30px);
}
.cropper-box{
  width: 100%;
  height: 500px;
}
</style>
