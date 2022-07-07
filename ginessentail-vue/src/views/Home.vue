<template>
  <div class="home">
    <div class="mydiv" @click="$router.replace({ name: 'javatosql' })">
      java to sql
    </div>

    <div class="mydiv" @click="$router.replace({ name: 'compareFile' })">
      对比文件
    </div>

    <div class="mydiv" @click="$router.replace({ name: 'daode' })">
      dao de jing
    </div>
    <div class="mydiv" @click="$router.replace({ name: 'cal' })">
      calculate
    </div>

  </div>
</template>

<script>
// @ is an alias to /src

export default {
  name: 'Home',
  components: {
  }
}
</script>

<style>
.mydiv {
  cursor: pointer;
  margin-top: 10px;
}
</style>
