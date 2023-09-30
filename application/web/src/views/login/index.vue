<template>
  <div class="login-container">
    <el-form ref="loginForm" class="login-form" auto-complete="on" label-position="left">

      <div class="title-container">
        <h3 class="title">欢迎使用区块链企业合同系统</h3>
      </div>
      

      <el-button :loading="loading" type="primary" style="width:100%;margin-bottom:30px;" @click.native.prevent="handleLogin">立即进入</el-button>

      <!-- <div class="tips">
        <span style="margin-right:20px;">tips: 选择不同用户角色模拟交易</span>
      </div> -->

    </el-form>
    <el-dialog title="欢迎使用" :visible="loginDialogVisible" @close="closeLoginDialog">
      <el-tabs v-model="activeTab">
        <el-tab-pane label="企业登录" name="company">
          <el-form :model="companyLoginForm" ref="companyLoginForm" label-width="80px">
            <el-form-item label="账号">
              <!-- <el-input placeholder="请输入账号"></el-input> -->
              <el-select v-model="value" placeholder="请选择已注册企业" class="login-select" @change="selectGet">
                <el-option
                  v-for="item in accountList"
                  :key="item.accountId"
                  :label="item.userName"
                  :value="item.accountId"
                >
                </el-option>
              </el-select>
            </el-form-item>
            <el-form-item label="密码">
              <el-input type="password" placeholder="请输入密码"></el-input>
            </el-form-item>
            <div class="dialog-footer">
              <el-button @click="showRegisterDialog">注册</el-button>
              <el-button type="primary" @click="companyLogin">登录</el-button>
            </div>
          </el-form>
        </el-tab-pane>
        <el-tab-pane label="部门登录" name="department">
          <el-form :model="departmentLoginForm" ref="departmentLoginForm" label-width="80px">
            <el-form-item label="账号">
              <el-input placeholder="请输入账号"></el-input>
            </el-form-item>
            <el-form-item label="密码">
              <el-input type="password" placeholder="请输入密码"></el-input>
            </el-form-item>
            <el-form-item label="企业码">
              <el-input placeholder="请输入对应企业的企业码"></el-input>
            </el-form-item>
            <div class="dialog-footer">
              <el-button @click="showRegisterDialog">注册</el-button>
              <el-button type="primary" @click="departmentLogin">登录</el-button>
            </div>
          </el-form>
        </el-tab-pane>
      </el-tabs>
    </el-dialog>
    <el-dialog title="企业/部门注册" :visible="registerDialogVisible" style="width: 1000px; align-self: center;" @close="closeRegisterDialog">
      <el-tabs v-model="activeTab">
        <el-tab-pane label="企业注册" name="company">
          <el-form :model="companyLoginForm" ref="companyLoginForm" label-width="80px">
            <el-form-item label="账号">
              <el-input placeholder="请输入账号"></el-input>
            </el-form-item>
            <el-form-item label="密码">
              <el-input type="password" placeholder="请输入密码"></el-input>
            </el-form-item>
            <div class="dialog-footer">
              <!-- <el-button @click="registerNewAccount">注册</el-button> -->
              <el-button type="primary" @click="registerNewCompany">注册</el-button>
            </div>
          </el-form>
        </el-tab-pane>
        <el-tab-pane label="部门注册" name="department">
          <el-form :model="departmentRegisterForm" ref="departmentRegisterForm" label-width="80px">
            <el-form-item label="账号">
              <el-input placeholder="请输入账号"></el-input>
            </el-form-item>
            <el-form-item label="密码">
              <el-input type="password" placeholder="请输入密码"></el-input>
            </el-form-item>
            <el-form-item label="企业码">
              <el-input placeholder="请输入对应企业的企业码"></el-input>
            </el-form-item>
            <div class="dialog-footer">
              <!-- <el-button @click="registerNewAccount">注册</el-button> -->
              <el-button type="primary" @click="registerNewDepartment">注册</el-button>
            </div>
          </el-form>
        </el-tab-pane>
      </el-tabs>
    </el-dialog>
  </div>
</template>

<script>
import { queryAccountList } from '@/api/account'

export default {
    name: 'Login',
    data() {
        return {
            loading: false,
            redirect: undefined,
            accountList: [],
            value: '',
            loginDialogVisible: false,
            registerDialogVisible: false,
            loginForm: {
                username: '',
                password: ''
            },
        };
    },
    watch: {
        $route: {
            handler: function (route) {
                this.redirect = route.query && route.query.redirect;
            },
            immediate: true
        }
    },
    created() {
        queryAccountList().then(response => {
            if (response !== null) {
                this.accountList = response;
            }
        });
    },
    methods: {
        handleLogin() {
            this.showLoginDialog();
        },
        selectGet(accountId) {
            this.value = accountId;
        },
        showLoginDialog() {
            this.loginDialogVisible = true;
        },
        closeLoginDialog() {
            this.loginDialogVisible = false;
        },
        companyLogin() {
            if (this.value) {
                this.loading = true;
                this.$store.dispatch('account/login', this.value).then(() => {
                    this.$router.push({ path: this.redirect || '/' });
                    this.loading = false;
                }).catch(() => {
                    this.loading = false;
                });
            }
            else {
                this.$message('请选择用户角色');
            }
        },
        showRegisterDialog() {
            this.registerDialogVisible = true;
            this.loginDialogVisible = false;
        },
        closeRegisterDialog() {
            this.registerDialogVisible = false;
        },
        registerNewCompany() {
          // registerNewAccount Logic
          closeRegisterDialog();
        },
        registerNewDepartment() {
          closeRegisterDialog();
        },
        departmentLogin() {
        }
    },
}
</script>

<style lang="scss" scoped>
$bg:#2d3a4b;
$dark_gray:#889aa4;
$light_gray:#eee;

.login-container {
  min-height: 100%;
  width: 100%;
  background-color: $bg;
  overflow: hidden;

  .login-form {
    position: relative;
    width: 520px;
    max-width: 100%;
    padding: 160px 35px 0;
    margin: 0 auto;
    overflow: hidden;
  }
  // .login-select{
  //  padding: 20px 0px 30px 0px;
  //  min-height: 100%;
  //  width: 100%;
  //  background-color: $bg;
  //  overflow: hidden;
  //  text-align: center;
  // }
  .tips {
    font-size: 14px;
    color: #fff;
    margin-bottom: 10px;

    span {
      &:first-of-type {
        margin-right: 16px;
      }
    }
  }

  .svg-container {
    padding: 6px 5px 6px 15px;
    color: $dark_gray;
    vertical-align: middle;
    width: 30px;
    display: inline-block;
  }

  .title-container {
    position: relative;

    .title {
      font-size: 26px;
      color: $light_gray;
      margin: 0px auto 40px auto;
      text-align: center;
      font-weight: bold;
    }
  }

  .show-pwd {
    position: absolute;
    right: 10px;
    top: 7px;
    font-size: 16px;
    color: $dark_gray;
    cursor: pointer;
    user-select: none;
  }
}
</style>
