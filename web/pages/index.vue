<template>
  <el-tabs>
    <el-tab-pane label="list">
      <el-input
        placeholder="Filter keyword"
        v-model="filterText">
      </el-input>

      <el-tree
        class="filter-tree"
        :props="defaultProps"
        :filter-node-method="filterNode"
        ref="tree"
        lazy
        :load="loadNode"
        @node-click="copy">
        <span slot-scope="{ node, data }">
          <span>{{ node.label }}</span>
        </span>
      </el-tree>
    </el-tab-pane>
    <el-tab-pane label="search">
      <el-form :inline="true" @submit.native.prevent>
        <el-form-item label="Sound Name">
          <el-input v-model="searchName" @keyup.enter.native="onSubmit"></el-input>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="onSubmit">Search</el-button>
        </el-form-item>
      </el-form>
      <div id="search-result">
        <span v-for="sound in searchResult" :key="sound.ID" class="search-node">
          <el-button
            type="text"
            @click="() => copy(sound, null, null, true)"
            class="search-node-button">
            {{ sound.name}}
          </el-button>
        </span>
      </div>
    </el-tab-pane>
  </el-tabs>
</template>

<script>
import axios from 'axios'
  export default {
    watch: {
      filterText(val) {
        this.$refs.tree.filter(val);
      }
    },

    methods: {
      filterNode(value, data) {
        if (!value) return true;
        return data.name.indexOf(value) !== -1;
      },
      loadNode(node, resolve){
        if (node.level === 0){
          axios.get(`${process.env.baseUrl}/categories/`)
            .then((res) => {
              resolve(res.data.result);
            });
        }else if(node.level === 1){
          axios.get(`${process.env.baseUrl}/categories/${node.data.ID}`)
            .then((res) => {
              const sounds = res.data.result.sounds.map((sound) => {
                sound.leaf = true;
                return sound;
              });
              resolve(sounds);
            });
        }
      },
      copy(data, node, treeNode, isSearch){
        if(isSearch || data.leaf){
          this.$copyText(data.name)
            .then((e) => {
              this.$message(`Copied ${data.name}`);
            }, (e) => {
              console.log(e);
            });
        }
      },
      async onSubmit(){
        let { data }  = await axios.get(`${process.env.baseUrl}/sounds/?query=${this.searchName}`);
        this.searchResult = data.result;
      }
    },

    data() {
      return {
        filterText: '',
        defaultProps: {
          children: 'sounds',
          label: 'name',
          isLeaf: 'leaf'
        },
        searchName: '',
        searchResult: []
      };
    }
  };
</script>

<style>
  .el-header {
    background-color: #B3C0D1;
    color: #333;
    line-height: 60px;
  }
  .search-node {
    flex: 1;
    display: flex;
    align-items: center;
    justify-content: space-between;
    line-height: 0%;
  }
  .search-node-button{
    line-height: 0;
  }
</style>
