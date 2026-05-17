using StackExchange.Redis;

public class RedisService
{
    private readonly IDatabase _db;
    // public IDatabase Db => _db;
    public RedisService(string connectionString)
    {
        var connection = ConnectionMultiplexer.Connect(connectionString);
        _db = connection.GetDatabase();
    }
    public void XAdd(string stream, string payload)
    {
        _db.StreamAdd(stream,new NameValueEntry[]
        {
            new NameValueEntry("payload", payload)
        });
    }
    public StreamEntry[] XRead(string stream, string lastId = "$")
    {
        var result =  _db.StreamRead(stream, lastId,count: 10);
        return result ?? Array.Empty<StreamEntry>();  
    }
}