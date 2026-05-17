using StackExchange.Redis;

public class MessageWorker : BackgroundService
{
    private readonly RedisService _redis;
    private readonly ILogger<MessageWorker> _logger;
    private string _lastId = "0";
    public MessageWorker(RedisService redis,ILogger<MessageWorker> logger)
    {
        _redis = redis;
        _logger = logger;
    }
    protected override async Task ExecuteAsync(CancellationToken stoppingToken)
    {
       _logger.LogInformation("Monitorando mensagens...");
       while (!stoppingToken.IsCancellationRequested)
        {
            var entries = _redis.XRead("wa:messages", _lastId);
        
            foreach(var entry in entries)
            {
                _lastId = entry.Id;

                var payload = entry.Values
                .FirstOrDefault(v => v.Name == "payload")
                .Value.ToString();
            
                _logger.LogInformation($"Processando mensagem: {payload}");
            }
            await Task.Delay(1000, stoppingToken); 
        }
    }
}