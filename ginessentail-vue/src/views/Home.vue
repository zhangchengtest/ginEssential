<template>
  <div class="home">
    <div class="mydiv" @click="$router.replace({ name: 'javatosql' })">
      java to sql
    </div>

    <div class="mydiv" @click="$router.replace({ name: 'compareFile' })">
      对比文件
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
}
</style>
