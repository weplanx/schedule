<?php

namespace lumen\extra\Redis;

use lumen\extra\common\RedisModel;

final class Sms extends RedisModel
{
    protected $key = 'sms:';

    /**
     * Factory Verify Code For Phone
     * @param string $phone PhoneNumber
     * @param string $code Code
     * @param int $timeout Timeout
     * @return bool
     */
    public function factory($phone, $code, $timeout = 120)
    {
        $data = json_encode([
            'code' => $code,
            'publish_time' => time(),
            'timeout' => $timeout
        ]);

        return $this->redis->setex($this->key . $phone, $timeout, $data);
    }

    /**
     * Verify Code
     * @param string $phone PhoneNumber
     * @param string $code Code
     * @param boolean $once Only Once
     * @return bool
     */
    public function check($phone, $code, $once = false)
    {
        if (!$this->redis->exists($this->key . $phone)) {
            return false;
        }

        $data = json_decode($this->redis->get($this->key . $phone), true);
        $result = ($code === $data['code']);
        if ($once && $result) {
            $this->redis->del($this->key . $phone);
        }

        return $result;
    }

    /**
     * Get Time Information
     * @param string $phone PhoneNumber
     * @return array|bool
     */
    public function time($phone)
    {
        if (!$this->redis->exists($this->key . $phone)) {
            return false;
        }

        $data = json_decode($this->redis->get($this->key . $phone), true);
        return [
            'publish_time' => $data['publish_time'],
            'timeout' => $data['timeout']
        ];
    }
}
