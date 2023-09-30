<template>
    <div class="contract-container">
      <el-row justify="center">
        <el-col :span="2999">
          <el-card>
            <div slot="header" class="clearfix">
              <span>合同信息</span>
            </div>
            <!-- 省略合同信息显示部分 -->
  
            <!-- 添加发起合同按钮 -->
            <div class="item" v-if="!editingContract">
              <el-button type="primary" @click="startEditingContract">发起合同</el-button>
            </div>
            <!-- 控制合同编辑界面的显示 -->
            <div v-if="editingContract">
              <el-button @click="cancelEditingContract">取消编辑</el-button>
              <el-dialog
                title="编辑合同"
                :visible="editingContract"
                @close="cancelEditingContract"
              >
                <!-- 在这里添加编辑合同的表单 -->
                <el-form :model="contractData" label-width="100px">
                  <el-form-item label="合同名称">
                    <el-input v-model="contractData.contractName"></el-input>
                  </el-form-item>
                  <el-form-item label="合同内容">
                    <el-input type="textarea" v-model="contractData.contractContent"></el-input>
                  </el-form-item>
                  <el-form-item label="创建者名称">
                    <el-input v-model="contractData.creatorName" disabled></el-input>
                  </el-form-item>
                  <el-form-item label="创建者公司">
                    <el-input v-model="contractData.creatorCompany" disabled></el-input>
                  </el-form-item>
                  <el-form-item label="创建时间">
                    <el-input v-model="contractData.createTime" disabled></el-input>
                  </el-form-item>
                  <el-form-item>
                    <el-button type="primary" @click="submitContract">提交合同</el-button>
                  </el-form-item>
                </el-form>
              </el-dialog>
            </div>
          </el-card>
        </el-col>
      </el-row>
    </div>
  </template>
  
  <script>
  export default {
    data() {
      return {
        editingContract: false,
        contractData: {
          contractName: '合同名称', // 合同名称，自主填写
          contractContent: '合同内容', // 合同内容，自主填写
          creatorName: '自动获取的创建者名称', // 自动获取的创建者名称
          creatorCompany: '自动获取的创建者公司', // 自动获取的创建者公司
          createTime: '自动获取的创建时间', // 自动获取的创建时间
        },
      };
    },
    methods: {
      getCurrentTime() {
        const now = new Date();
        const year = now.getFullYear();
        const month = String(now.getMonth() + 1).padStart(2,'0');
        const day = String(now.getDate()).padStart(2, '0');
        const hours = String(now.getHours()).padStart(2,'0');
        const minutes = String(now.getMinutes()).padStart(2,'0');
        const seconds = String(now.getSeconds()).padStart(2,'0');
        this.contractData.createTime = `${year}-${month}-${day} ${hours}:${minutes}:${seconds}`;
      },
      startEditingContract() {
        this.getCurrentTime();
        this.editingContract = true;
      },
      cancelEditingContract() {
        this.editingContract = false;
      },
      // 可以添加其他处理编辑合同的方法
    },
  };
  </script>
  
  <style scoped>
  .contract-container {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100vh;
  }
  
  .item {
    margin-bottom: 10px;
  }
  </style>
  