<template>
  <div class="config-edit">
    <ElForm label-position="top" :model="editConfig">
      <ElFormItem label="账户">
        <ElInput v-model="editConfig.username"></ElInput>
      </ElFormItem>
      <ElFormItem label="管理地址">
        <ElInput v-model="editConfig.web"></ElInput>
      </ElFormItem>
      <ElFormItem label="抓包间隔">
        <ElInput v-model="editConfig.capture_interval" type="number" style="width: 80px;"></ElInput><span>s</span>
      </ElFormItem>
      <ElFormItem label="检测间隔">
        <ElInput v-model="editConfig.check_interval" type="number" style="width: 80px;"></ElInput><span>s</span>
      </ElFormItem>
      <ElFormItem label="是否修改密码">
        <ElCheckbox v-model="modifyPassword">是</ElCheckbox>
      </ElFormItem>
      <ElFormItem v-if="modifyPassword" label="旧密码">
        <ElInput :type="oldPasswordVisible ? 'text' : 'password'" v-model="oldPassword">
          <template #append>
            <ElButton @click="toggleOldPasswordVisibility" type="primary" circle>
              <el-icon v-if="oldPasswordVisible"><View /></el-icon>
            </ElButton>
          </template>
        </ElInput>
      </ElFormItem>
      <ElFormItem v-if="modifyPassword" label="新密码">
        <ElInput :type="newPasswordVisible ? 'text' : 'password'" v-model="newPassword">
          <template #append>
            <ElButton @click="toggleNewPasswordVisibility" type="primary" circle>
              <el-icon v-if="newPasswordVisible"><View /></el-icon>
            </ElButton>
          </template>
        </ElInput>
      </ElFormItem>
      <ElFormItem v-if="modifyPassword" label="确认新密码">
        <ElInput :type="confirmNewPasswordVisible ? 'text' : 'password'" v-model="confirmNewPassword">
          <template #append>
            <ElButton @click="toggleConfirmNewPasswordVisibility" type="primary" circle>
              <el-icon v-if="confirmNewPasswordVisible"><View /></el-icon>
            </ElButton>
          </template>
        </ElInput>
      </ElFormItem>
      <ElFormItem label="黑名单">
        <ElTable :data="localBlacklist.map(ip => ({ ip }))" style="width: 100%">
          <ElTableColumn prop="ip" label="IP" width="180"/>
          <ElTableColumn fixed="right" label="操作" width="100">
            <template #default="scope">
              <ElButton @click="removeBlacklistItem(scope.$index)" type="danger" size="small">删除</ElButton>
            </template>
          </ElTableColumn>
        </ElTable>
        <ElInput v-model="newBlacklistItem" placeholder="添加 IP 到黑名单">
          <template #append>
            <ElButton @click="addBlacklistItem" type="primary">添加</ElButton>
          </template>
        </ElInput>
      </ElFormItem>
      <ElFormItem>
        <ElButton type="primary" @click="confirmEdit">提交</ElButton>
      </ElFormItem>
    </ElForm>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import axios from 'axios';
import { ElForm, ElFormItem, ElInput, ElButton, ElMessage, ElMessageBox, ElTable, ElTableColumn, ElCheckbox } from 'element-plus';
import { View } from '@element-plus/icons-vue';
import { user } from '../utils';

interface Config {
  username: string;
  capture_interval: string;
  check_interval: string;
  web: string;
  tix: Array<{ type: string; blacklist?: string[]; apikey?: string }>;
}

const router = useRouter();

const editConfig = ref<Config>({
  username: '',
  capture_interval: '',
  check_interval: '',
  web: '',
  tix: []
});

const localBlacklist = ref<string[]>([]);

const oldPassword = ref('');
const newPassword = ref('');
const confirmNewPassword = ref('');

const oldPasswordVisible = ref(false);
const newPasswordVisible = ref(false);
const confirmNewPasswordVisible = ref(false);

const modifyPassword = ref(false);

const newBlacklistItem = ref('');

const toggleOldPasswordVisibility = () => {
  oldPasswordVisible.value = !oldPasswordVisible.value;
};

const toggleNewPasswordVisibility = () => {
  newPasswordVisible.value = !newPasswordVisible.value;
};

const toggleConfirmNewPasswordVisibility = () => {
  confirmNewPasswordVisible.value = !confirmNewPasswordVisible.value;
};

const addBlacklistItem = () => {
  if (newBlacklistItem.value && !localBlacklist.value.includes(newBlacklistItem.value)) {
    localBlacklist.value.push(newBlacklistItem.value);
    newBlacklistItem.value = '';
  }
};

const removeBlacklistItem = (index: number) => {
  localBlacklist.value.splice(index, 1);
};

onMounted(() => {
  axios.get('/api/config', {
    headers: {
      Authorization: `Bearer ${user.value.token}`
    }
  }).then(res => {
    const config = res.data;
    const localConfig = config.tix.find((t: any) => t.type === 'local');
    editConfig.value = {
      username: config.username,
      capture_interval: config.capture_interval,
      check_interval: config.check_interval,
      web: config.web,
      tix: config.tix
    };
    if (localConfig && localConfig.blacklist) {
      localBlacklist.value = localConfig.blacklist;
    }
  }).catch(e => {
    if (e.response.status === 401) {
      router.push('/login');
    } else {
      ElMessage.error(e.response.data.error);
    }
  });
});

const confirmEdit = () => {
  if (modifyPassword.value && (newPassword.value !== confirmNewPassword.value)) {
    ElMessage.error('两次密码输入不匹配');
    return;
  }

  ElMessageBox.confirm(
    '是否提交修改？',
    '提示',
    {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning',
    }
  ).then(() => {
    submitEdit();
  }).catch(() => {
    ElMessage.info('已取消修改');
  });
};

const submitEdit = () => {
  const updates: Partial<Config> & { oldPassword?: string; newPassword?: string } = {
    username: editConfig.value.username,
    capture_interval: editConfig.value.capture_interval,
    check_interval: editConfig.value.check_interval,
    web: editConfig.value.web,
    tix: editConfig.value.tix.map(t => t.type === 'local' ? { ...t, blacklist: localBlacklist.value } : t)
  };

  if (modifyPassword.value) {
    updates.oldPassword = oldPassword.value;
    updates.newPassword = newPassword.value;
  }

  axios.post('/api/config', updates, {
    headers: {
      Authorization: `Bearer ${user.value.token}`
    }
  }).then(() => {
    ElMessage.success('配置已更新');
    router.push('/tix');
  }).catch(e => {
    ElMessage.error(e.response.data.error);
  });
};
</script>

<style scoped>
.config-edit {
  display: flex;
  flex-direction: column;
  gap: 20px;
}
</style>
