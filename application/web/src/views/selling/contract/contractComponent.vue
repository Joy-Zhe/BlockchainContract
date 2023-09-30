<template>
    <div class="contract-container">
      <el-row justify="center">
        <el-col :span="3000">
          <el-card>
            <div slot="header" class="clearfix">
              <span>合同信息</span>
            </div>
            <div class="item">
              <el-tag type="success">合同甲方: </el-tag>
              <span>{{ partyA }}</span>
            </div>
            <div class="item">
              <el-tag type="warning">合同乙方: </el-tag>
              <span>{{ partyB }}</span>
            </div>
            <div class="item">
              <el-tag type="info">状态: </el-tag>
              <span>{{ contractStatus }}</span>
            </div>
            <div class="item" v-if="contractStatus === '待签订'">
              <el-button type="primary" @click="agreeContract">同意签订</el-button>
              <el-button type="danger" @click="rejectContract">拒绝签订</el-button>
              <el-button @click="openModifyDialog">修改合同</el-button>
            </div>
            <div class="item" v-else>
              <el-button @click="openModifyDialog">修改合同</el-button>
            </div>
          </el-card>
        </el-col>
      </el-row>
      <el-dialog
        title="修改合同"
        :visible="modifyDialogVisible"
        @close="modifyDialogVisible = false"
      >
        <!-- 在这里添加修改合同的表单或内容 -->
        <!-- 例如： -->
        <!-- <el-form :model="formData">
          ... 表单内容 ...
        </el-form> -->
      </el-dialog>
    </div>
  </template>
  
  <script>
  export default {
    data() {
      return {
        partyA: '甲方公司',
        partyB: '乙方公司',
        contractStatus: '待签订', // 初始状态为待签订，可根据实际情况修改
        modifyDialogVisible: false, // 控制修改合同对话框的显示状态
        // formData: {}, // 修改合同的表单数据，可以在这里定义
        currentTime: '',
      };
    },
    methods: {
      agreeContract() {
        // 同意签订合同的逻辑
        this.contractStatus = '已签订';
      },
      rejectContract() {
        // 拒绝签订合同的逻辑
        this.contractStatus = '已拒绝';
      },
      openModifyDialog() {
        // 打开修改合同对话框
        this.modifyDialogVisible = true;
      },
      // 可以添加其他处理修改合同的方法
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
  