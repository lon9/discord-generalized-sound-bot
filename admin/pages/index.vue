<template>
  <el-container>
    <el-header style="text-align:right;">
      <el-button class="right" @click="logout">Logout</el-button>
    </el-header>
    <el-main>
      <el-form ref="form" :model="form" label-width="120px">
        <el-form-item label="Category">
          <el-select
            v-model="form.categoryName"
            filterable
            allow-create
            default-first-option
            remote
            clearable
            no-data-text="No category"
            placeholder="Category">
            <el-option
              v-for="category in categories"
              :key="category.ID"
              :label="category.name"
              :value="category.name">
            </el-option>
          </el-select>
        </el-form-item>
        <el-form-item label="Sound name">
          <el-input
            placeholder="Sound name"
            v-model="form.name"
            clearable>
          </el-input>
        </el-form-item>
        <el-form-item label="Sound file">
          <el-upload
            ref="upload"
            drag
            action="/adimin/sounds"
            :auto-upload="false"
            accept="audio/*">
            <i class="el-icon-upload"></i>
            <div class="el-upload__text">Drop file here or <em>click to upload</em></div>
            <div class="el-upload__tip" slot="tip">wav/mp3/ogg/flac file</div>
          </el-upload>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="onSubmit">Send</el-button>
        </el-form-item>
      </el-form>
    </el-main>
  </el-container>
</template>

<script>
import axios from 'axios'
export default {
  methods: {
    onSubmit(){
      if(this.$refs.upload.uploadFiles.length == 0) return;
      let formData = new FormData();
      formData.append('file', this.$refs.upload.uploadFiles[0].raw);
      formData.append('name', this.form.name);
      formData.append('categoryName', this.form.categoryName);
      axios.post(`${process.env.baseUrl}/admin/sounds`, formData, {headers: {Authorization: 'Bearer ' + this.$store.state.token}})
        .then((res) => {
          this.$message(`Uploaded ${this.form.name}`);
          this.form.name = '';
          this.form.categoryName = '';
          this.$refs.upload.clearFiles();
        }).catch((err) => {
          console.error(err);
        });
    },
    async logout(){
      try{
        await this.$store.dispatch('logout')
        this.$nuxt.$router.replace({ path: '/login'})
      }catch(e){
        console.error(e);
      }
    }
  },
  data () {
    return {
      form: {
        name: '',
        categoryName: ''
      } 
    }
  },
  async asyncData(){
    let {data} = await axios.get(`${process.env.baseUrl}/categories/`);
    return {categories: data.result};
  },
  middleware: 'auth'
}
</script>

<style>
  .el-header {
    background-color: #B3C0D1;
    color: #333;
    line-height: 60px;
  }
</style>
