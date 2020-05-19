const Koa = require('koa');
const router = require('koa-router')();
const app = new Koa();
const routes = require('./routes/router.js');
const cors = require('koa2-cors')
const bodyParser = require('koa-bodyparser');
//中间件（匹配路由之前操作）
app.use(async (ctx,next)=>{
  console.log("中间件")
  await next();

  if (ctx.status == 404){
    ctx.body = "<h1>页面未找到</h1>"
  }else{
    console.log(ctx.url);
  }
})

app.use(bodyParser());
//跨域
app.use(cors());
//配置路由

router.use('/', routes);

//启动路由

app.use(router.routes());
app.use(router.allowedMethods());


app.listen(3000,()=>{
  console.log('start-quick is starting at port 3000')
})
