var builder = WebApplication.CreateBuilder(args);

builder.Services.AddControllers();
builder.Services.AddSingleton<RedisService>(_ => new RedisService("localhost:6379"));
builder.Services.AddHostedService<MessageWorker>();

var app = builder.Build();

app.MapControllers();

app.Run();