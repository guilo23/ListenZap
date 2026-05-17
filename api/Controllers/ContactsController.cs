using Microsoft.AspNetCore.Mvc;
using System.Text.Json;

[ApiController]
[Route("api/[controller]")]
public class ContactsController : ControllerBase
{
    private readonly RedisService _redis;
    public ContactsController(RedisService redis)
    {
        _redis = redis;
    }
    [HttpPost]
    public IActionResult Post([FromBody] Contacts contact)
    {
        var payload = JsonSerializer.Serialize(contact);
        _redis.XAdd("wa:contacts", payload);
        return Accepted();
    }
}